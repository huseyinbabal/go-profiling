import http from 'k6/http';

export const options = { vus: 4, duration: '20s', summaryTrendStats: ['avg', 'p(95)'] };

export default function () {
  http.get('http://localhost:9092/search?q=ERROR.*took=4[0-9]{2}ms');
}
