import express from "express";
import axios from "axios";
import dotenv from "dotenv";
import bodyParser from "body-parser";
import dayjs from "dayjs";

dotenv.config();
const app = express();
app.use(bodyParser.json());

app.get("/testapi", async (req, res) => {
    res.send("api is working");
});

const {
    MPESA_CONSUMER_KEY,
    MPESA_CONSUMER_SECRET,
    MPESA_SHORTCODE,
    MPESA_PASSKEY,
    MPESA_ENV,
    CALLBACK_BASE_URL,
    GO_BACKEND_NOTIFY_URL,
    NODE_NOTIFY_SECRET,
    JWT_SECRET,
} = process.env;

const MPESA_BASE =
    MPESA_ENV === "production"
        ? "https://api.safaricom.co.ke"
        : "https://sandbox.safaricom.co.ke";

// auth token
async function getAccessToken() {
    const url = `${MPESA_BASE}/oauth/v1/generate?grant_type=client_credentials`;
    const auth = Buffer.from(`${MPESA_CONSUMER_KEY}:${MPESA_CONSUMER_SECRET}`).toString("base64");
    const res = await axios.get(url, {
        headers: { Authorization: `Basic ${auth}` },
    });
    return res.data.access_token;
}

// Generate Lipa Na Mpesa password
function lipaPassword() {
    const timestamp = dayjs().format("YYYYMMDDHHmmss");
    const raw = `${MPESA_SHORTCODE}${MPESA_PASSKEY}${timestamp}`;
    return { password: Buffer.from(raw).toString("base64"), timestamp };
}

// STK Push Initiation
app.post("/mpesa/stkpush", async (req, res) => {
    try {
        const { order_id, amount, phone } = req.body;
        if (!order_id || !amount || !phone) {
            return res.status(400).json({ error: "order_id, amount and phone are required" });
        }

        const token = await getAccessToken();
        const { password, timestamp } = lipaPassword();
        const callbackUrl = `${CALLBACK_BASE_URL}/mpesa/callback`;

        const payload = {
            BusinessShortCode: MPESA_SHORTCODE,
            Password: password,
            Timestamp: timestamp,
            TransactionType: "CustomerPayBillOnline",
            Amount: amount,
            PartyA: phone,
            PartyB: MPESA_SHORTCODE,
            PhoneNumber: phone,
            CallBackURL: callbackUrl,
            AccountReference: order_id,
            TransactionDesc: `Payment for order ${order_id}`,
        };

        const resp = await axios.post(
            `${MPESA_BASE}/mpesa/stkpush/v1/processrequest`,
            payload,
            { headers: { Authorization: `Bearer ${token}` } }
        );

        return res.json({ daraja: resp.data });
    } catch (err) {
        console.error("stkpush error:", err.response?.data || err.message);
        return res.status(500).json({
            error: "failed to initiate stkpush",
            details: err.response?.data || err.message,
        });
    }
});

// Fixed callback handler
app.post("/mpesa/callback", async (req, res) => {
    res.json({ ResultCode: 0, ResultDesc: "Accepted" });

    try {
        const body = req.body;
        console.log("MPESA CALLBACK:", JSON.stringify(body, null, 2));

        const stkCallback = body?.Body?.stkCallback;
        if (!stkCallback) {
            console.warn("missing stkCallback");
            return;
        }

        const checkoutRequestID = stkCallback.CheckoutRequestID;
        const merchantRequestID = stkCallback.MerchantRequestID;
        const resultCode = stkCallback.ResultCode;
        const status = resultCode === 0 ? "success" : "failed";

        const items = stkCallback.CallbackMetadata?.Item || [];
        const mpesaReceipt = items.find(i => i.Name === "MpesaReceiptNumber")?.Value;
        const phone = items.find(i => i.Name === "PhoneNumber")?.Value;
        const accountRef = items.find(i => i.Name === "AccountReference")?.Value;
        const amount = items.find(i => i.Name === "Amount")?.Value;

        console.log("Extracted from callback:", {
            accountRef,
            checkoutRequestID,
            merchantRequestID,
            amount,
            phone,
            mpesaReceipt
        });

        // The real order ID should be in AccountReference
        if (!accountRef) {
            console.error("AccountReference is missing from M-Pesa callback - this contains the order_id");
            console.error("Available items:", items);
            return;
        }

        const orderId = accountRef; // This should be your UUID order ID

        // Validate that orderId is a proper UUID
        const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
        if (!uuidRegex.test(orderId)) {
            console.error(`Invalid order ID format: ${orderId}. Expected UUID format.`);
            return;
        }

        // Fetch order details using the correct order ID
        let orderTotal = amount; // Use M-Pesa amount as fallback
        try {
            const orderResp = await axios.get(
                `${GO_BACKEND_NOTIFY_URL.replace('/payments/confirm', '')}/orders/${orderId}`,
                {
                    headers: {
                        Authorization: `Bearer ${JWT_SECRET}`,
                        "X-Node-Notify-Secret": NODE_NOTIFY_SECRET,
                    }
                }
            );

            orderTotal = orderResp.data.order?.total || amount;
            console.log(`Fetched order total: ${orderTotal} for order: ${orderId}`);
        } catch (err) {
            console.error("Failed to fetch order details from Go backend:", err.response?.data || err.message);
            console.log(`Using M-Pesa amount as fallback: ${amount}`);
        }

        console.log("Notifying Go backend with:", {
            order_id: orderId,
            status,
            amount: orderTotal,
            provider: "mpesa"
        });

        // Notify backend
        const notifyResponse = await axios.post(
            GO_BACKEND_NOTIFY_URL,
            {
                order_id: orderId,
                status,
                amount: orderTotal,
                provider: "mpesa",
                checkout_request_id: checkoutRequestID,
                merchant_request_id: merchantRequestID,
                mpesa_receipt: mpesaReceipt,
                phone: String(phone),
                raw: body
            },
            {
                headers: {
                    "X-Node-Notify-Secret": NODE_NOTIFY_SECRET,
                    "Content-Type": "application/json"
                }
            }
        );

        console.log("Go backend notification successful:", notifyResponse.status);

    } catch (err) {
        console.error("Error processing mpesa callback:", err.response?.data || err.message);
        if (err.response) {
            console.error("Response status:", err.response.status);
            console.error("Response data:", err.response.data);
        }
    }
});

const port = process.env.PORT || 5000;
app.listen(port, () => console.log(`mpesa service running on port ${port}`));
