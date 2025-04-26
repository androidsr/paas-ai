export default function chartPlugin(md) {
    md.block.ruler.before('fence', 'chart', chartBlock, {
        alt: ['paragraph', 'reference', 'blockquote', 'list']
    });

    md.renderer.rules.chart = (tokens, idx) => {
        const token = tokens[idx];
        const codeId = token.meta.codeId;
        const chartData = encodeURIComponent(token.content);
        return `<div class="chart-container" id="${codeId}" data-chart-data="${chartData}"></div>`;
    };
}

function chartBlock(state, startLine, endLine, silent) {
    const pos = state.bMarks[startLine] + state.tShift[startLine];
    const max = state.eMarks[startLine];

    const marker = '```chart';
    const lineStart = state.src.slice(pos, max).trim();
    if (!lineStart.startsWith(marker)) return false;

    let nextLine = startLine + 1;
    let endLineIndex = -1;

    while (nextLine < endLine) {
        const lineText = state.src.slice(state.bMarks[nextLine], state.eMarks[nextLine]).trim();
        if (lineText === '```') {
            endLineIndex = nextLine;
            break;
        }
        nextLine++;
    }

    if (endLineIndex === -1) return false;
    if (silent) return true;

    const content = state.getLines(startLine + 1, nextLine, state.tShift[startLine], true);
    const codeId = toShortBase36Hash(content);

    const token = state.push('chart', 'div', 0);
    token.block = true;
    token.content = content.trim();
    token.map = [startLine, nextLine + 1];
    token.meta = { codeId };

    state.line = nextLine + 1;
    return true;
}

function toShortBase36Hash(str) {
    const hash = Array.from(str).reduce((acc, c) => acc + c.charCodeAt(0), 0);
    return 'chart-' + (hash + str.length).toString(36);
}

