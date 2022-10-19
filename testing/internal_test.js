import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
    stages: [
      // Ramp-up from 1 to 5 virtual users (VUs) in 5s
      { duration: "5s", target: 20 },
  
      // Stay at rest on 5 VUs for 10s
      { duration: "10s", target: 5 },
  
      // Ramp-down from 5 to 0 VUs for 5s
      { duration: "5s", target: 0 },
    ],
  };
  
  const payload = JSON.stringify(
    {
      "method": "sendSignedTxn",
      "parameter": ["deez nuts signed"]
      }
  )

  const url = "http://localhost:5000/v1/transaction/rcvSignedTxn"

  export default function () {
    const response = http.post(url, payload, {
      headers: { Accepts: "application/json" },
    });
    check(response, { "status is 200": (r) => r.status === 200 });
    sleep(0.3);
  }