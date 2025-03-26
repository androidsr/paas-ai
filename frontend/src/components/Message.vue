<template>
    <div class="chat-container">
        <div class="chat-box" ref="chatBox">
            <div v-for="(item, index) in messages" :key="index" class="message-container" :class="item.type">
                <UserOutlined v-if="item.type == 'user'" style="padding-left: 10px;font-size: 20px;color: cadetblue;" />
                <AndroidOutlined v-if="item.type == 'ai'"
                    style="padding-right: 10px;font-size: 20px;color: cadetblue;" />
                <div class="message-content" v-if="item.type == 'user' || item.type == 'ai'">
                    <div v-html="renderMarkdown(item.message)"></div>
                    <div v-if="item.type == 'ai'"><a class="copy-message" @click="copyMessage(index)"
                            v-if="item.message.length > 0">复制</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import Clipboard from 'clipboard';
import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark.css"; // 更贴近 VSCode 主题
import MarkdownIt from "markdown-it";
import mila from "markdown-it-link-attributes";

export default {
    props: {
        messages: {
            type: Array,
            required: true,
        },
    },
    mounted() {
        if (this.copyCode == null) {
            this.copyCode = new Clipboard('.copy-btn')
            this.copyCode.on('success', (e) => {
                this.$message.success('复制成功')
            })
            this.copyCode.on('error', (e) => {
                this.$message.error('复制成功失败')
            });
        }
    },
    destroyed() {
        this.copyCode.destroy()
    },
    data() {
        return {
            copyCode: null,
            mdParser: new MarkdownIt({
                html: false,
                linkify: true,
                breaks: true,
                xhtmlOut: true,
                typographer: true,
                highlight: function (str, lang) {
                    const codeIndex = `copy-${Date.now()}-${Math.floor(Math.random() * 10000000)}`;

                    let copyButton = `<a class="copy-btn" style="float:right; cursor:pointer;" 
                        data-clipboard-action="copy" data-clipboard-target="#${codeIndex}">复制</a>\n`;

                    let highlightedCode = str;

                    if (lang && hljs.getLanguage(lang)) {
                        try {
                            highlightedCode = hljs.highlight(str, { language: lang }).value;
                        } catch (error) {
                            console.error("代码高亮错误:", error);
                        }
                    }
                    return `
                        <pre class="hljs" style="padding:8px !important">
                            <code>${copyButton}<span id="${codeIndex}">${highlightedCode}</span></code>
                        </pre>
                    `;
                }
            }).use(mila, { attrs: { target: "_blank", rel: "noopener" } }),
        };
    },
    methods: {
        copyMessage(index) {
            let clipboard = new Clipboard('.copy-message', {
                text: () => this.messages[index].message || "",
            });
            clipboard.on('success', (e) => {
                this.$message.success('复制成功');
                clipboard.destroy();

            })
            clipboard.on('error', (e) => {
                this.$message.error('复制成功失败');
                clipboard.destroy();
            });
        },
        renderMarkdown(content) {
            return this.mdParser.render(content);
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
    padding: 10px 10px;
    border-radius: 12px;
    word-wrap: break-word;
    position: relative;
    font-size: 14px;
    line-height: 1.6;
}

.user .message-content {
    background: white;
    color: black;
    text-align: right;
    border-top-right-radius: 0;
    border: 1px solid #d9d5d5;

}

.ai .message-content {
    background: white;
    color: black;
    border-top-left-radius: 0;
    border: 1px solid #d9d5d5;
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

pre.hljs {
    padding: 12px 2px 12px 40px !important;
    border-radius: 5px !important;
    position: relative;
    font-size: 14px !important;
    line-height: 22px !important;
    overflow: hidden !important;

    code {
        display: block !important;
        margin: 0 10px !important;
        overflow-x: auto !important;
    }

    .line-numbers-rows {
        position: absolute;
        pointer-events: none;
        top: 12px;
        bottom: 12px;
        left: 0;
        font-size: 100%;
        width: 40px;
        text-align: center;
        letter-spacing: -1px;
        border-right: 1px solid rgba(0, 0, 0, .66);
        user-select: none;
        counter-reset: linenumber;

        span {
            pointer-events: none;
            display: block;
            counter-increment: linenumber;

            &:before {
                content: counter(linenumber);
                color: #999;
                display: block;
                text-align: center;
            }
        }
    }

    b.name {
        position: absolute;
        top: 2px;
        right: 50px;
        z-index: 10;
        color: #999;
        pointer-events: none;
    }

    .copy-btn {
        position: absolute;
        top: 2px;
        right: 4px;
        z-index: 10;
        color: #333;
        cursor: pointer;
        background-color: #fff;
        border: 0;
        border-radius: 2px;
    }
}
</style>
