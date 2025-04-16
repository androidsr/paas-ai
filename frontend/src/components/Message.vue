<template>
    <div class="chat-container">
        <div class="chat-box" ref="chatBox">
            <div v-for="(item, index) in formattedMessages" :key="item.id" class="message-container" :class="item.type">
                <UserOutlined v-if="item.type === 'user'" :style="userIconStyle" />
                <AndroidOutlined v-if="item.type === 'ai'" :style="aiIconStyle" />
                <div class="message-content" ref="messageContainer">
                    <MdContent :renderedMessage="item.renderedMessage" :key="'content_' + item.id"
                        :ref="'content_' + item.id"></MdContent>
                    <div v-if="item.type === 'ai'">
                        <a class="copy-message" @click="copyMessage(item.message)" v-if="item.message.length > 0">复制</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import Chart from "@/components/Chart.vue";
import MermaidChat from "@/components/MermaidChat.vue";

import MdContent from "@/components/MdContent.vue";
import Clipboard from 'clipboard';
import hljs from "highlight.js";
import "highlight.js/styles/atom-one-dark.css";
import MarkdownIt from "markdown-it";
import mila from "markdown-it-link-attributes";
import { createVNode, defineComponent, h, render } from 'vue';

export default defineComponent({
    components: {
        Chart,
        MermaidChat,
        MdContent,
    },
    props: {
        messages: {
            type: Array,
            required: true,
        },
    },
    computed: {
        userIconStyle() {
            return { paddingLeft: '10px', fontSize: '20px', color: 'cadetblue' };
        },
        aiIconStyle() {
            return { paddingRight: '10px', fontSize: '20px', color: 'cadetblue' };
        },
        formattedMessages() {
            return this.messages.map((item) => ({
                ...item,
                renderedMessage: this.renderMarkdown(item.message),
            }));
        },
    },
    mounted() {
        this.initClipboard();
        this.$nextTick(() => {
            this.renderAll();
        });
    },
    updated() {
        this.$nextTick(() => {
            this.renderVueComponents();
        });
    },
    beforeUnmount() {
        if (this.copyCode) {
            this.copyCode.destroy();
        }
    },
    data() {
        return {
            copyCode: null,
            isRendering: false,
            mdParser: new MarkdownIt({
                html: true,
                linkify: true,
                breaks: true,
                xhtmlOut: true,
                typographer: true,
                highlight: this.highlightCode,
            }).use(mila, { attrs: { target: "_blank", rel: "noopener" } }),
        };
    },
    methods: {
        initClipboard() {
            if (this.copyCode) {
                this.copyCode.destroy();
            }
            this.copyCode = new Clipboard('.copy-btn');
            this.copyCode.on('success', () => {
                this.$message.success('复制成功');
            });

            this.copyCode.on('error', () => {
                this.$message.error('复制失败');
            });
        },
        copyMessage(message) {
            const copy = new Clipboard('.copy-message', {
                text: () => {
                    return message;
                },
            });
            copy.on('success', () => {
                this.$message.success('复制成功');
                copy.destroy();
            });

            copy.on('error', () => {
                this.$message.error('复制失败');
                copy.destroy();
            });
        },
        renderMarkdown(content) {
            return content ? this.mdParser.render(content) : "";
        },

        highlightCode(str, lang) {
            if (lang === "chart") {
                try {
                    const data = JSON.parse(str);
                    const chartId = `chart-${this.messages.length}`;
                    return `<div class="chart-container" data-chart-id="${chartId}" data-chart-data='${JSON.stringify(data)}'></div>`;
                } catch (error) {
                    return '...';
                }
            } else if (lang === 'mermaid') {
                const chartId = `mermaid-${this.messages.length}`;
                const data = encodeURIComponent(str);
                return `<div class="mermaid-container" data-chart-id="${chartId}" data-chart-data="${data}"></div>`;
            }

            let codeId = `code-${Date.now()}-${Math.floor(Math.random() * 10000)}`;
            let copyButton = `<a class="copy-btn" style="float:right; cursor:pointer;" data-clipboard-action="copy" data-clipboard-target="#${codeId}">复制</a>`;

            let highlightedCode = str;
            if (lang && hljs.getLanguage(lang)) {
                try {
                    highlightedCode = hljs.highlight(str, { language: lang }).value;
                } catch (error) {
                    console.error("代码高亮错误:", error);
                }
            }

            return `<pre class="hljs" style="padding:2px !important">${copyButton}<code id="${codeId}">${highlightedCode}</code></pre>`;
        },
        renderAll() {
            this.isRendering = false;
            this.renderVueComponents();
        },
        capitalizeFirstLetter(str) {
            if (typeof str !== 'string' || str.length === 0) {
                return str;
            }
            return str.charAt(0).toUpperCase() + str.slice(1);
        },
        renderVueComponents() {
            if (this.isRendering) return; // 如果正在渲染中，跳过这次渲染
            this.isRendering = true;

            this.$nextTick(() => {
                if (!this.$refs.messageContainer) {
                    this.isRendering = false;
                    return;
                }

                const messageContainers = Array.isArray(this.$refs.messageContainer)
                    ? this.$refs.messageContainer
                    : [this.$refs.messageContainer];

                messageContainers.forEach((container) => {
                    container.querySelectorAll('.chart-container').forEach((el) => {
                        const chartData = JSON.parse(el.dataset.chartData);
                        const chartId = el.dataset.chartId;///"chart-" + this.messages.length;
                        const existingChart = this.$refs[chartId];
                        if (existingChart && existingChart.chartId == chartId) {
                            return; // 数据未变化，不重新渲染
                        }
                        const chartContainer = document.createElement('div');
                        el.replaceWith(chartContainer);

                        // 更新 chartData
                        this.$nextTick(() => {
                            const chartVNode = h(Chart, {
                                chartData,
                                key: chartId,  // 使用 chartId 作为 key 来保证更新时唯一性
                                ref: chartId,  // 给每个图表元素加一个 ref 用来缓存
                            });

                            const vm = createVNode(chartVNode); // 创建虚拟节点
                            const chartElement = document.createElement('div');
                            vm.appContext = this.$.appContext; // 保证上下文一致
                            render(vm, chartElement); // 使用 render 渲染虚拟节点
                            chartContainer.appendChild(chartElement);
                        });
                    });

                    container.querySelectorAll('.mermaid-container').forEach((el) => {
                        const chartData = el.dataset.chartData;
                        const chartId = el.dataset.chartId;
                        const existingChart = this.$refs[chartId];
                        if (existingChart && existingChart.chartId == chartId) {
                            return;
                        }
                        if (!chartData) {
                            return;
                        }
                        const chartContainer = document.createElement('div');
                        el.replaceWith(chartContainer);

                        // 更新 chartData
                        this.$nextTick(() => {
                            const chartVNode = h(MermaidChat, {
                                chartData,
                                key: chartId,  // 使用 chartId 作为 key 来保证更新时唯一性
                                ref: chartId,  // 给每个图表元素加一个 ref 用来缓存
                            });

                            const vm = createVNode(chartVNode); // 创建虚拟节点
                            const chartElement = document.createElement('div');
                            vm.appContext = this.$.appContext; // 保证上下文一致
                            render(vm, chartElement); // 使用 render 渲染虚拟节点
                            chartContainer.appendChild(chartElement);
                        });
                    });
                });
                this.isRendering = false;
            });
        },
    },
});
</script>


<style scoped>
.chat-container {
    display: flex;
    flex-direction: column;
    overflow: hidden;
    height: 77vh;
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

.message-content {
    max-width: 94%;
    padding: 10px;
    border-radius: 12px;
    word-wrap: break-word;
    position: relative;
    font-size: 14px;
    line-height: 1.6;
    box-shadow: 1px 2px 1px rgba(73, 164, 234, 0.3);
}

.user .message-content {
    background: white;
    color: black;
    text-align: right;
    border-top-right-radius: 0;
    border: 1px solid #93b1f1;
}

.ai .message-content {
    min-width: 70%;
    background: white;
    color: black;
    border-top-left-radius: 0;
    border: 1px solid #93b1f1;
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
    padding: 12px;
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
</style>
