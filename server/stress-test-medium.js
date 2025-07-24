import http from "k6/http";
import { check, sleep } from "k6";

// Medium load test - next level after gentle test succeeds
export let options = {
  stages: [
    { duration: "1m", target: 20 }, // Ramp to 20 users
    { duration: "3m", target: 40 }, // Peak at 40 users
    { duration: "2m", target: 20 }, // Scale down
    { duration: "1m", target: 0 }, // Ramp down
  ],
  thresholds: {
    http_req_duration: ["p(95)<3000"], // 95% under 3s
    http_req_failed: ["rate<0.05"], // Error rate under 5%
  },
};

const BASE_URL = "http://localhost:4300";

export default function () {
  // Health check
  let healthRes = http.get(`${BASE_URL}/api/v1/health`);
  check(healthRes, {
    "health check status is 200": (r) => r.status === 200,
  });

  sleep(0.3);

  // Articles
  let articlesRes = http.get(`${BASE_URL}/api/v1/articles`);
  check(articlesRes, {
    "articles status is 200": (r) => r.status === 200,
  });

  sleep(0.5);

  // Categories
  let categoriesRes = http.get(`${BASE_URL}/api/v1/categories`);
  check(categoriesRes, {
    "categories status is 200": (r) => r.status === 200,
  });

  sleep(0.3);

  // Tags
  let tagsRes = http.get(`${BASE_URL}/api/v1/tags`);
  check(tagsRes, {
    "tags status is 200": (r) => r.status === 200,
  });

  // Shorter random delay
  sleep(Math.random() * 1);
}
