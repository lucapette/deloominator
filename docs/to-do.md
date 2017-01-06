# To Do

## Issues to create

Here is a list of things I want to fix but I can do out of the alpha phase as
they are non critical to the product.

- All the test with postgres and mysql work only with default credentials. A
  good solution is to rely on the data sources coming from the env.
- Make the mysql connector work with different protocols (the DSN works
  differently there)
- Add an automate step to the build process in order to have a full list of
  env vars the app understands. It's good to have it that way so we can
  validate the input at startup better.
