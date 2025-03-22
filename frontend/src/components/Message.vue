<template>
    <div class="chat-container">
        <div class="chat-box" ref="chatBox">
            <div v-for="(msg, index) in messages" :key="index" class="message-container" :class="msgClass(index)">
                <!-- <img :src="avatars[msgClass(index)]" class="avatar" /> -->
                <UserOutlined v-if="msgClass(index) == 'user'"
                    style="padding-left: 10px;font-size: 22px;color: cadetblue;" />
                <AndroidOutlined v-if="msgClass(index) == 'ai'"
                    style="padding-right: 10px;font-size: 22px;color: cadetblue;" />
                <div class="message-content">
                    <div v-html="renderMarkdown(msg)"></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import MarkdownIt from "markdown-it";
import hljs from "highlight.js";
import mila from "markdown-it-link-attributes";
import "highlight.js/styles/atom-one-dark.css"; // 更贴近 VSCode 主题

export default {
    props: {
        messages: {
            type: Array,
            required: true,
        },
    },
    data() {
        return {
            mdParser: new MarkdownIt({
                html: false,
                linkify: true,
                breaks: true,
                xhtmlOut: true,
                typographer: true,
                highlight: (str, lang) => {
                    if (lang && hljs.getLanguage(lang)) {
                        return `<pre class="hljs" style='padding:12px'><code>${hljs.highlight(str, { language: lang }).value}</code></pre>`;
                    }
                    return `<pre class="hljs" style='padding:12px'><code>${hljs.highlightAuto(str).value}</code></pre>`;
                },
            }).use(mila, { attrs: { target: "_blank", rel: "noopener" } }),
        };
    },
    methods: {
        renderMarkdown(content) {
            return this.mdParser.render(content);
        },
        msgClass(index) {
            return index % 2 === 0 ? "user" : "ai";
        },
        scrollToBottom() {
            this.$nextTick(() => {
                const chatBox = this.$refs.chatBox;
                chatBox.scrollTop = chatBox.scrollHeight;
            });
        },
    },
    watch: {
        messages() {
            this.scrollToBottom();
        },
    },
};
</script>

<style scoped>
.chat-container {
    display: flex;
    flex-direction: column;
    overflow: hidden;
    height: 78vh;
}

.chat-box {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
}

.message-container {
    display: flex;
    align-items: flex-start;
    margin-bottom: 12px;
}

.user {
    flex-direction: row-reverse;
}

.avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin: 0 10px;
}

.message-content {
    max-width: 90%;
    padding: 10px 14px;
    border-radius: 12px;
    word-wrap: break-word;
    position: relative;
    font-size: 14px;
    line-height: 1.6;
}

.user .message-content {
    background: #007aff;
    color: white;
    text-align: right;
    border-top-right-radius: 0;
}

.ai .message-content {
    background: #f1f1f1;
    color: black;
    border-top-left-radius: 0;
}

.hljs {
    background: #1e1e1e !important;
    color: #d4d4d4;
    border-radius: 8px;
    overflow-x: auto;
    font-family: "Fira Code", "Consolas", monospace;
    font-size: 14px;
    line-height: 1.6;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    border: 1px solid #444;
    /* 边框 */
    padding: 40px;
    /* 这里是 hljs 的 padding */
}

.hljs code {
    white-space: pre-wrap;
    /* 使代码块可以换行 */
    word-wrap: break-word;
}
</style>
