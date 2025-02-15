<template>
    <div class="markdown-container" ref="markdownContainer" v-show="!!rawContent && parsedContent">
        <div v-html="parsedContent" class="markdown-content"></div>
    </div>
    <!-- 复制按钮 -->
    <!-- <a-button type="text" @click="copyAllText" icon="copy" class="copy-all-button">复制</a-button> -->
</template>

<script>
import { message } from 'ant-design-vue';
import hljs from 'highlight.js'; // 引入 highlight.js
import { marked } from 'marked';
import 'highlight.js/styles/atom-one-dark.css'; // 你可以根据需求选择不同的主题样式

export default {
    props: {
        content: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            parsedContent: '', // 存储渲染后的 HTML 内容
            rawContent: this.$store.chats.result, // 存储原始的 Markdown 内容
        };
    },
    watch: {
        content(newContent) {
            this.appendMarkdown(newContent);
        },
    },
    mounted() {
        this.parseMarkdown(this.content);
    },
    methods: {
        parseMarkdown(content) {
            let parsed = marked(content);
            this.parsedContent = parsed;

            // 使用 $nextTick 确保 DOM 更新完成后再进行高亮
            this.$nextTick(() => {
                this.highlightCode();
            });
        },

        appendMarkdown(newContent) {
            this.rawContent += newContent; // 累加新的原始内容
            this.$store.setResult(this.rawContent);  // 更新 Vuex 存储中的内容
            this.parseMarkdown(this.rawContent);    // 重新解析累加后的内容
        },

        copyAllText() {
            const textToCopy = this.$refs.markdownContainer.innerText;
            navigator.clipboard.writeText(textToCopy)
                .then(() => {
                    message.success("复制成功！");
                })
                .catch((err) => {
                    message.error("复制失败！");
                });
        },

        highlightCode() {
            const codeBlocks = this.$refs.markdownContainer.querySelectorAll('pre code');
            codeBlocks.forEach((block) => {
                hljs.highlightElement(block); // 高亮每个代码块
            });
        },

        // 清空内容方法
        clearContent() {
            this.rawContent = '';   // 清空原始内容
            this.parsedContent = ''; // 清空渲染后的内容
            this.$store.setResult(''); // 如果您使用 Vuex 存储内容，清空存储中的内容
        },
    },
};
</script>

<style scoped>
.markdown-container {
    padding: 8px;
    border-radius: 2px;
    border: 1px solid #ccc;
    overflow-y: auto;
    max-height: 87vh;
}

.markdown-content {
    line-height: 1.6;
    font-family: 'Arial', sans-serif;
    color: #333;
}

.copy-all-button {
    margin-top: 10px;
}

pre {
    background-color: #2d2d2d;
    color: #f8f8f2;
    padding: 20px;
    border-radius: 5px;
    overflow-x: auto;
    white-space: pre-wrap;  /* 确保长代码换行 */
    word-wrap: break-word;   /* 允许换行 */
}

code {
    font-family: 'Courier New', monospace;
    display: block;
    padding: 0 10px;
}

h1,
h2,
h3,
h4,
h5,
h6 {
    font-weight: bold;
}

ul,
ol {
    padding-left: 20px;
}

li {
    margin-bottom: 8px;
}

blockquote {
    margin: 10px 0;
    padding: 10px;
    background-color: #f9f9f9;
    border-left: 5px solid #ccc;
}
</style>
