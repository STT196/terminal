import { SESv2Client, SendEmailCommand } from "@aws-sdk/client-sesv2";

export namespace Email {
  export async function send(
    template: string,
    to: string,
    subject: string,
    body: string,
  ) {
    const client = new SESv2Client();
    await client.send(
      new SendEmailCommand({
        FromEmailAddress: "hello@terminal.shop",
        Destination: {
          ToAddresses: [to],
        },
        Content: {
          Simple: {
            Subject: { Data: subject },
            Body: { Text: { Data: body } },
          },
        },
      }),
    );
  }
}
