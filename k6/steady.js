import http from 'k6/http';

export const options = {
  summaryTrendStats: ['avg', 'p(95)'],
  scenarios: {
    steady: { executor: 'constant-arrival-rate', rate: 10, timeUnit: '1s', duration: '60s', preAllocatedVUs: 10 },
  },
};

export default function () {
  http.get('http://localhost:9092/search?q=WARN');
}
