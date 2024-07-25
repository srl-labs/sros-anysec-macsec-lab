#!/bin/sh

# Start /app/server and redirect its logs
/app/server 2>&1 | tee /var/log/server.log &

# Start frontend and redirect its logs
npm run preview 2>&1 | tee /var/log/frontend.log &

# Wait for all background processes to finish
# effectively blocks the CMD from finishing
wait