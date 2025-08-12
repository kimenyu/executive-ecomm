import express from "express";
import axios from "axios";
import dotenv from "dotenv";
import bodyParser from "body-parser";
import dayjs from "dayjs";
import jwt from "jsonwebtoken";

dotenv.config();
const app = express();
app.use(bodyParser.json());

// In-memory store for mapping CheckoutRequestID to order_id
// In production, use Redis or database
const checkoutOrderMap = new Map();

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

// Generate JWT token for Go backend authentication
function generateJWTToken() {
    const payload = {
        service: "mpesa-node",
        iat: Math.floor(Date.now() / 1000),
        exp: Math.floor(Date.now() / 1000) + (60 * 60) // 1 hour expiration
    };
    return jwt.sign(payload, JWT_SECRET);
}
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
            AccountReference: order_id, // Still send order_id, but don't rely on it in callback
            TransactionDesc: `Payment for order ${order_id}`,
        };

        console.log("STK Push payload:", JSON.stringify(payload, null, 2));

        const resp = await axios.post(
            `${MPESA_BASE}/mpesa/stkpush/v1/processrequest`,
            payload,
            { headers: { Authorization: `Bearer ${token}` } }
        );

        // Store the mapping between CheckoutRequestID and order_id
        const checkoutRequestID = resp.data.CheckoutRequestID;
        if (checkoutRequestID) {
            checkoutOrderMap.set(checkoutRequestID, {
                order_id,
                amount,
                phone,
                timestamp: new Date().toISOString()
            });
            console.log(`Stored mapping: ${checkoutRequestID} -> ${order_id}`);
        }

        return res.json({ daraja: resp.data });
    } catch (err) {
        console.error("stkpush error:", err.response?.data || err.message);
        return res.status(500).json({
            error: "failed to initiate stkpush",
            details: err.response?.data || err.message,
        });
    }
});

// Callback handler
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

        // Get order details from our mapping
        const orderMapping = checkoutOrderMap.get(checkoutRequestID);
        if (!orderMapping) {
            console.error(`No order mapping found for CheckoutRequestID: ${checkoutRequestID}`);
            console.error("Available mappings:", Array.from(checkoutOrderMap.keys()));
            return;
        }

        const { order_id: orderId, amount: originalAmount } = orderMapping;
        console.log(`Found order mapping: ${checkoutRequestID} -> ${orderId}`);

        // Extract callback data
        const items = stkCallback.CallbackMetadata?.Item || [];
        const mpesaReceipt = items.find(i => i.Name === "MpesaReceiptNumber")?.Value;
        const phone = items.find(i => i.Name === "PhoneNumber")?.Value;
        const callbackAmount = items.find(i => i.Name === "Amount")?.Value;

        console.log("Extracted from callback:", {
            orderId,
            checkoutRequestID,
            merchantRequestID,
            callbackAmount,
            originalAmount,
            phone,
            mpesaReceipt,
            status
        });

        // Verify amount matches
        if (callbackAmount && parseFloat(callbackAmount) !== parseFloat(originalAmount)) {
            console.warn(`Amount mismatch: callback=${callbackAmount}, original=${originalAmount}`);
        }

        // Use the original amount from our mapping since we have it
        // Skip fetching from Go backend to avoid JWT complexity
        const orderTotal = originalAmount;
        console.log(`Using original amount from mapping: ${orderTotal}`);

        // Notify Go backend
        const notifyPayload = {
            order_id: orderId,
            status,
            amount: orderTotal,
            provider: "mpesa",
            checkout_request_id: checkoutRequestID,
            merchant_request_id: merchantRequestID,
            mpesa_receipt: mpesaReceipt,
            phone: String(phone),
            raw: body
        };

        console.log("Notifying Go backend with:", JSON.stringify(notifyPayload, null, 2));

        const notifyResponse = await axios.post(
            GO_BACKEND_NOTIFY_URL,
            notifyPayload,
            {
                headers: {
                    "X-Node-Notify-Secret": NODE_NOTIFY_SECRET,
                    "Content-Type": "application/json"
                }
            }
        );

        console.log("Go backend notification successful:", notifyResponse.status);

        // Clean up the mapping after successful processing
        checkoutOrderMap.delete(checkoutRequestID);
        console.log(`Cleaned up mapping for: ${checkoutRequestID}`);

    } catch (err) {
        console.error("Error processing mpesa callback:", err.response?.data || err.message);
        if (err.response) {
            console.error("Response status:", err.response.status);
            console.error("Response headers:", err.response.headers);
        }
    }
});

// Debug endpoint to check stored mappings
app.get("/mpesa/mappings", (req, res) => {
    const mappings = Array.from(checkoutOrderMap.entries()).map(([checkoutId, data]) => ({
        checkoutRequestID: checkoutId,
        ...data
    }));
    res.json({ mappings, count: mappings.length });
});

// Cleanup old mappings (run every hour)
setInterval(() => {
    const now = new Date();
    let cleaned = 0;

    for (const [checkoutId, data] of checkoutOrderMap.entries()) {
        const age = now - new Date(data.timestamp);
        // Remove mappings older than 1 hour (3600000 ms)
        if (age > 3600000) {
            checkoutOrderMap.delete(checkoutId);
            cleaned++;
        }
    }

    if (cleaned > 0) {
        console.log(`Cleaned up ${cleaned} old checkout mappings`);
    }
}, 3600000);

const port = process.env.PORT || 5000;
app.listen(port, () => console.log(`mpesa service running on port ${port}`));