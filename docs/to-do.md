# To Do

## Issues to create

Here is a list of things I want to fix but I can do out of the alpha phase as
they are non critical to the product.

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

## Investigate and change

- Having the app in the context for the graphql handler feels wrong
- Not sure why we need a type switch in the Query method and why mysql and
  postgres get different structures from the database
