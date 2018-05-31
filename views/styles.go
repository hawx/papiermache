package views

const Styles = `
body {
    font: 16px/1.3em Georgia;
    margin: 4rem;
}

header {
  margin: 1rem 0 2rem;
}

nav a {
    color: black;
    text-decoration: none;
    font-style: italic;
}

nav a + a {
    margin-left: .5rem;
}

nav .selected {
    border-bottom: 1px solid #ccc;
}

ul {
    list-style: none;
    padding: 0;
    margin: 3rem 0;
    max-width: 40rem;
}

li {
    position: relative;
    display: flex;
    border-top: 1px solid #ccc;
}

li:last-child {
    border-bottom: 1px solid #ccc;
}

li > a {
    color: black;
    text-decoration: none;
    padding: 1rem 0;
    display: block;
    word-break: keep-all;
    text-overflow: ellipsis;
    overflow-x: hidden;
    margin: 0 1rem;
}

.actions {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    margin-right: 1rem;
    text-align: right;
}
`
