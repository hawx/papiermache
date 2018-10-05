package views

const Styles = `
body {
    font: 16px/1.3em Georgia;
    color: #0a0708;
    margin: 0;
    background: #eee;
}

body > header {
    box-sizing: border-box;
    margin: 3rem auto;
    padding: 0 1rem;
    max-width: 40rem;
}

body > header + ul {
    padding-top: 0;
}

body > header {
    background: #eee;
}

body > header nav {
    display: flex;
    justify-content: space-between;
}

header h1 a {
    color: black;
    text-decoration: none;
}

nav a {
    color: black;
    text-decoration: none;
    font-style: italic;
}

nav a + a {
    margin-left: .5rem;
}

nav span {
    margin-left: .5rem;
    margin-right: .5rem;
    color: #aaa;
}

nav .selected {
    border-bottom: 1px solid #aaa;
}

ul {
    list-style: none;
    padding: 0;
    margin: 0 auto;
    max-width: 40rem;
}

li {
    position: relative;
    display: flex;
    border: 1px solid #ccc;
    background: #fafafa;
    margin-bottom: .5rem;
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
    margin-left: 1rem;
}

.actions {
    display: flex;
    margin-left: .5rem;
    width: 3.5rem;
}

.actions a {
    margin: 1rem .5rem;
    display: inline-block;
}

.actions svg {
    height: 1rem;
    width: 1rem;
    opacity: .5;
}

.actions a:hover svg {
    opacity: 1;
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
}


/* article */

article {
    background: #fafafa;
    border: 1px solid #ccc;
    margin: 3rem auto;
    padding: 2rem 3rem;
    max-width: 60rem;
    font-size: 21px;
    line-height: 1.4;
}

article header {
    position: relative;
}

article header h1 + .meta {
    margin-top: -1rem;
    margin-bottom: 2rem;
}

article header .meta span {
    color: #333;
}

article header a {
    color: #333;
}

article header time {
    color: #333;
    font-style: italic;
}

article header .actions {
    display: inline-flex;
}

.content {

}

article img {
    max-width: 60rem;
}

`
