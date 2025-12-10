const Cerebras = require('@cerebras/cerebras_cloud_sdk');

const cerebras = new Cerebras({
    apiKey: "csk-2rdnnrej2wf3v34rhcxyk5v5hje9cv2yfrr8f3pr5hk2r294"
});

async function test() {
    try {
        console.log("Testing Cerebras...");
        const response = await cerebras.chat.completions.create({
            messages: [{role: "user", content: "Say hello!"}],
            model: 'llama-3.3-70b',
            max_tokens: 50
        });
        console.log("Success:", response.choices[0].message.content);
    } catch (e) {
        console.error("Error:", e.message);
    }
}

test();
