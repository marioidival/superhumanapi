import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 50,
  duration: '60s',
};


export default function() {
  const emails = [
    'contatodaraingrid@gmail.com',
    'marioidival@gmail.com',
    'javalisson@gmail.com',
    'rahul@superhuman.com',
    'conrad@superhuman.com',
  ];

  const email = emails[Math.floor(Math.random() * emails.length)];
  http.get('http://localhost:3000/lookup?email=' + email);
  sleep(1);
}
