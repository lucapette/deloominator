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
  validate the input at startup better. We can use the Usage function for
  envconfing (it's pretty nice!)
- testutil.InitApp has a strange signature. Look for improvements there.
- `// Format() does not work both ways yet` this comment makes it clear we
  need to find a better way of handling the differences between databases DSNs
- Add Stringer interface to Rows Row and Column

## Missing documents

- Explain how to use `/bin/run` script in `docs/developer-manual.md`
