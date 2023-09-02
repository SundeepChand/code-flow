import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';

mermaid.initialize({ startOnLoad: false });

const drawMermaidDiagram = async function(graphDefinition) {
    const { svg } = await mermaid.render('graphDiv', graphDefinition)
    document.getElementById('result-display').innerHTML = svg
}


const go = new Go();
WebAssembly.instantiateStreaming(fetch("assets/wasm/main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
});

window.addEventListener('load', function() {
    const btnSubmit = document.getElementById('btn-submit')

    console.log('Web page loaded')

    btnSubmit.addEventListener('click', async function() {
        const inputText = document.getElementById('code-input').value
        const result = generateMermaidCode(inputText)
        await drawMermaidDiagram(result)
    })
})