import random
import random
import sys
import time

import aiohttp


async def send_request(session, base_url):
    start_time = time.time()

    location_id = random.randint(1, 60)
    microcategory_id = random.randint(1, 60)
    user_id = random.randint(1, 300)

    params = {
        'location_id': location_id,
        'microcategory_id': microcategory_id,
        'user_id': user_id
    }

    headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }

    async with session.get(base_url, params=params, headers=headers) as response:
        response_text = await response.text()
        end_time = time.time()
        waiting_time = end_time - start_time
        print(f"Response: {response_text}, Waiting Time: {waiting_time} seconds")
        return waiting_time

async def main():
    if len(sys.argv) != 4:
        print("Usage: python script.py [server_port] [rps] [duration]")
        return

    base_url = f'http://localhost:{sys.argv[1]}/retrieve'
    rps = int(sys.argv[2])  # Requests per second
    duration = int(sys.argv[3])  # Duration in seconds
    num_requests = rps * duration  # Total number of requests

    waiting_times = []
    async with aiohttp.ClientSession() as session:
        start_time = time.time()
        for _ in range(num_requests):
            waiting_time = await send_request(session, base_url)
            waiting_times.append(waiting_time)
            await asyncio.sleep(1 / rps)  # Delay between requests to achieve the desired RPS

        end_time = time.time()
        print(f"Time taken: {end_time - start_time} seconds")
        actual_rps = num_requests / (end_time - start_time)
        print(f"Actual RPS: {actual_rps}")
        mean_waiting_time = sum(waiting_times) / len(waiting_times)
        print(f"Mean Waiting Time: {mean_waiting_time} seconds")

if __name__ == '__main__':
    asyncio.run(main())
