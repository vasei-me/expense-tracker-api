// در فایل extension.js یا main.js
const vscode = require("vscode");
const Cerebras = require("@cerebras/cerebras_cloud_sdk");

async function activate(context) {
  // دریافت API Key از تنظیمات یا محیط
  const config = vscode.workspace.getConfiguration("cerebras");
  const apiKey = config.get("apiKey") || process.env.CEREBRAS_API_KEY;

  const cerebras = new Cerebras({
    apiKey: apiKey,
  });

  // تعریف یک command
  let disposable = vscode.commands.registerCommand(
    "extension.askCerebras",
    async function () {
      // دریافت متن از کاربر
      const question = await vscode.window.showInputBox({
        prompt: "سوال خود را بپرسید",
      });

      if (question) {
        const completion = await cerebras.chat.completions.create({
          messages: [{ role: "user", content: question }],
          model: "llama-3.3-70b",
          max_completion_tokens: 1024,
          temperature: 0.2,
          top_p: 1,
          stream: false,
        });

        vscode.window.showInformationMessage(
          completion.choices[0].message.content
        );
      }
    }
  );

  context.subscriptions.push(disposable);
}
