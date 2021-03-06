const fs = require("fs")
const path = require("path")
fis.match("**", {
    release: false,
},1)
fis.match("(**).source.md", {
    parser: [
        function (content, file) {
            return content.replace(/\[(.*?)\|embed\]\((.*)\)/g, function (source, name, ref) {
                const extname = path.extname(ref)
                const code = fs.readFileSync(path.join(file.dirname, ref)).toString()
                return `[${name}](${ref})` +
                    '\r\n' +
                    '```' + extname +
                    '\r\n' +
                    code  +
                    '\r\n' +
                    '```'
            })
        }
    ],
    release: "/$1.md",
}, 999)