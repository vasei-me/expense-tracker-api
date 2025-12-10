// test-cerebras-esm.js
import Cerebras from "@cerebras/cerebras_cloud_sdk";

const cerebras = new Cerebras({
  apiKey: "csk-2rdnnrej2wf3v34rhcxyk5v5hje9cv2yfrr8f3pr5hk2r294",
});

async function main() {
  try {
    console.log("🔗 Connecting to Cerebras...");

    const completion = await cerebras.chat.completions.create({
      messages: [
        {
          role: "user",
          content:
            "Why is fast inference important for AI models? Answer in 2 sentences.",
        },
      ],
      model: "llama-3.3-70b",
      max_completion_tokens: 200,
      temperature: 0.2,
      top_p: 1,
      stream: false,
    });

    console.log("\n✅ Response received:\n");
    console.log(completion.choices[0].message.content);
  } catch (error) {
    console.error("❌ Error:", error.message);
    if (error.response) {
      console.error("Response status:", error.response.status);
      console.error("Response data:", error.response.data);
    }
  }
}

main();
