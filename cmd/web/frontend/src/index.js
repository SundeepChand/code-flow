import mermaid from "mermaid";
import { editor } from 'monaco-editor'

import './main.css'

// Initialise mermaid
mermaid.initialize({
    startOnLoad: false,
    theme: 'dark',
});

const drawMermaidDiagram = async function(graphDefinition) {
    const { svg, bindFunctions } = await mermaid.render('graphDiv', graphDefinition)
    bindFunctions
    document.getElementById('result-display').innerHTML = svg
}

// Load GO WASM binary
const go = new Go();
WebAssembly.instantiateStreaming(fetch("assets/wasm/main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

window.addEventListener('load', () => {
    // Hide the load screen
    document.getElementById('load-screen').style.display = 'none'

    // Setup the code editor
    const codeEditor = editor.create(document.getElementById('code-editor'), {
        value: `// HOW IT WORKS?
        // Write a go function (with package the name & import statement)
        // And click on the generate button. This will generate a flow chart of the given program
        // Currently only if-else and for statements are supported.
        // Support for other flow control statements (eg. switch, break, continue, etc.) coming soon.`,
        language: 'go',
        autoClosingBrackets: 'always',
        automaticLayout: true,
        padding: {
            top: 5,
            bottom: 0,
        },
        theme: 'vs-dark',
    })

    // Wire the event listeners
    const btnSubmit = document.getElementById('btn-submit')

    btnSubmit.addEventListener('click', async () => {
        const inputText = codeEditor.getValue()
        const result = generateMermaidCode(inputText)
        await drawMermaidDiagram(result)
    })
})
