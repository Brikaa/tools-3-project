cp .env.sample .env
if [ -f .env.overrides ]; then
  cat .env.overrides >> .env
fi
