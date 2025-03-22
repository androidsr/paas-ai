<template>
    <div class="chat-container">
        <div class="chat-box" ref="chatBox">
            <div v-for="(msg, index) in messages" :key="index" class="message-container" :class="msgClass(index)">
                <!-- <img :src="avatars[msgClass(index)]" class="avatar" /> -->
                <UserOutlined v-if="msgClass(index) == 'user'"
                    style="padding-left: 10px;font-size: 20px;color: cadetblue;" />
                <AndroidOutlined v-if="msgClass(index) == 'ai'"
                    style="padding-right: 10px;font-size: 20px;color: cadetblue;" />
                <div class="message-content">
                    <div v-html="renderMarkdown(msg)"></div>
                    <div><a class="copy-message" @click="copyMessage(index)" v-if="msg.length > 0">复制</a>
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
                    // 当前时间加随机数生成唯一的id标识
                    const codeIndex = parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
                    // 复制功能主要使用的是 clipboard.js
                    let html = `<a class="copy-btn" style="float:right" data-clipboard-action="copy" data-clipboard-target="#copy${codeIndex}">复制</a>\n`
                    if (lang && hljs.getLanguage(lang)) {
                        try {
                            // highlight.js 高亮代码
                            const preCode = hljs.highlight(lang, str, true).value
                            html = html + preCode
                            // 将代码包裹在 textarea 中
                            return `<pre class="hljs" style="padding:8px"><code>${html}</code></pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
                        } catch (error) {
                            console.log(error)
                        }
                    }

                    html = html + str;
                    // 将代码包裹在 textarea 中
                    return `<pre class="hljs" style="padding:8px"><code>${html}</code></pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(/<\/textarea>/g, '&lt;/textarea>')}</textarea>`
                }
            }).use(mila, { attrs: { target: "_blank", rel: "noopener" } }),
        };
    },
    methods: {
        copyMessage(index) {
            let clipboard = new Clipboard('.copy-message', {
                text: () => this.messages[index] || "",
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
    /* watch: {
        messages() {
           // this.scrollToBottom();
        },
    }, */
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
