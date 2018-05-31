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

#cover {
    top: 0;
    left: 0;
    z-index: 1000;
    position: absolute;
    height: 100%;
    width: 100%;
    background: rgba(0, 255, 255, .7);
    display: block;
    padding: 0;
    margin: 0;
}

#cover a {
    font: 16px/1.3em Helvetica, sans-serif;
    position: relative;
    display: block;
    left: 50%;
    top: 50%;
    text-align: center;
    width: 100px;
    margin-left: -50px;
    line-height: 50px;
    margin-top: -25px;
    font-size: 16px;
    font-weight: bold;
    border: 1px solid;
}`
