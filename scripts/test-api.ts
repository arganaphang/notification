import { faker } from "@faker-js/faker";

type Notification = {
  title: string;
  content: string;
  user_id: string;
  order_id: number;
};

const N = 100;
const BASE_URL = "http://127.1:8000";

async function main() {
  let all: Promise<Response>[] = [];
  for (let idx = 0; idx < N; idx++) {
    const title = faker.lorem.words(4);

    const data: Notification = {
      title: title.charAt(0).toUpperCase() + title.slice(1),
      content: faker.lorem.sentences(),
      user_id: faker.string.uuid(),
      order_id: faker.number.int(100),
    };

    all.push(
      fetch(`${BASE_URL}/notification`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
    );
  }

  try {
    await Promise.all(all);
    console.log("SUCCESS SEED");
  } catch (e) {
    console.log(e);
  }
}

main();
