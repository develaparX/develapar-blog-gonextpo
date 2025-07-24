import http from "k6/http";
import { check, sleep } from "k6";

// Gentle stress test configuration for limited VPS
export let options = {
  stages: [
    // Ramp up slowly
    { duration: "30s", target: 5 }, // Start with 5 users
    { duration: "1m", target: 10 }, // Increase to 10 users
    { duration: "2m", target: 15 }, // Peak at 15 users
    { duration: "1m", target: 10 }, // Scale down
    { duration: "30s", target: 0 }, // Ramp down
  ],
  thresholds: {
    http_req_duration: ["p(95)<2000"], // 95% of requests under 2s
    http_req_failed: ["rate<0.1"], // Error rate under 10%
  },
};

const BASE_URL = "http://localhost:4300";

export default function () {
  // Test different endpoints with delays

  // Health check (should be fast)
  let healthRes = http.get(`${BASE_URL}/api/v1/health`);
  check(healthRes, {
    "health check status is 200": (r) => r.status === 200,
  });

  sleep(0.5); // Small delay between requests

  // Get articles (main endpoint)
  let articlesRes = http.get(`${BASE_URL}/api/v1/articles`);
  check(articlesRes, {
    "articles status is 200": (r) => r.status === 200,
  });

  sleep(1); // Longer delay to reduce load

  // Get categories
  let categoriesRes = http.get(`${BASE_URL}/api/v1/categories`);
  check(categoriesRes, {
    "categories status is 200": (r) => r.status === 200,
  });

  sleep(0.5);

  // Random delay to simulate real user behavior
  sleep(Math.random() * 2);
}
