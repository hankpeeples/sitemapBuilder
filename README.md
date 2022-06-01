# Sitemap Builder

This program builds a sitemap of the website you give it. It visits the root page and makes a list of every link found on that page that points to a page on the same domain.
Once the list of links is created, it visits each of them and adds any new links it finds to the list.
This step repeats over and over therefore visiting every page on that domain that can be reached from the root page. (Given that the `-depth=?` flag allows the program to go deep enough to see pages).

### Usage
- Run the program using `go run main.go <optional flags> > map.xml`
<br>
The `> map.xml` will send the program output to a new file named `map.xml`. If this is not included, the xml format will be printed in the terminal window.
<br>
<br>
- The optional flags are as follows:
  - `-url=<your domain of choice>` The default is `https://pkg.go.dev`.
  - `-depth=<num>` The default is 1.
###
**Note:** Giving a depth of 2+ can cause a long run time depending on the domain you choose.
Using the default domain and a depth of 2 produced a total run time of about 3 minutes (791 links were visited during this time).