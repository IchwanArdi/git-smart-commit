// Pastikan Anda sudah mengatur process.env.GEMINI_API_KEY di lingkungan Anda
const apiKey = process.env.GEMINI_API_KEY;
const url = `https://generativelanguage.googleapis.com/v1beta/models/gemini-3.5-flash:generateContent`;

const headers = {
    "x-goog-api-key": apiKey,
    "Content-Type": "application/json"
};

const body = JSON.stringify({
    contents: [
        {
            parts: [
                {
                    text: "Say hello in one short sentence."
                }
            ]
        }
    ]
});

// Fungsi untuk mengeksekusi request
async function testGemini() {
    try {
        const response = await fetch(url, {
            method: "POST",
            headers: headers,
            body: body
        });

        const data = await response.json();
        console.log("Respon API:", JSON.stringify(data, null, 2));
    } catch (error) {
        console.error("Terjadi error:", error);
    }
}

testGemini();
