<template>
    <div>
        <div class="mermaid-container" @click="openModal">
            <div class="mermaid" v-html="svgContent" />
        </div>
        <a-modal v-model:open="showModal" title="Mermaid 预览" width="90%" :footer="null" centered @cancel="closeModal">
            <div class="zoom-wrapper">
                <div ref="zoomContainer" class="zoom-inner" v-html="svgContent" />
            </div>
        </a-modal>
    </div>
</template>

<script>
import mermaid from 'mermaid';
import panzoom from 'panzoom';

export default {
    name: 'MermaidChat',
    props: {
        chartData: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            svgContent: '',
            showModal: false,
            panZoomInstance: null,
            data: "",
        };
    },
    mounted() {
        this.renderMermaid();
    },
    watch: {
        chartData() {
            this.renderMermaid();
        },
    },
    methods: {
        changeChart() {
            this.$store.setLoading(true);
            let message = "这是一个mermaid图表的配置数据,请给我更新一种展现方式,并使用md代码块输出，语言类型为：mermaid。\n不要输出额外的信息。 图表数据：" + this.data;
            this.$bus.emit("changeChart", message);
            setTimeout(() => {
                this.$store.setLoading(false);
            }, 3000);
        },
        openModal() {
            this.showModal = true;
            this.$nextTick(() => {
                const container = this.$refs.zoomContainer;
                if (container) {
                    // 初始化缩放
                    this.panZoomInstance = panzoom(container, {
                        zoomSpeed: 0.05,
                        maxZoom: 5,
                        minZoom: 0.2,
                    });
                }
            });
        },
        closeModal() {
            this.showModal = false;
            if (this.panZoomInstance) {
                this.panZoomInstance.dispose();
                this.panZoomInstance = null;
            }
        },
        async renderMermaid() {
            if (!this.chartData) return;

            try {
                mermaid.initialize({ startOnLoad: false, theme: 'default' });
                this.data = decodeURIComponent(this.chartData);
                const { svg } = await mermaid.render(`mermaid-${Date.now()}`, this.data);
                this.svgContent = svg;
            } catch (e) {
                //console.error('Mermaid 渲染失败:', e);
                //this.svgContent = `<pre style="color: red;">Mermaid 渲染失败：${e.message}</pre>`;
            }
        },
    },

};
</script>

<style scoped>
.mermaid-container {
    width: 100%;
    overflow-x: auto;
    border: 1px solid #eee;
    background: #fdfdfd;
    border-radius: 6px;
    cursor: pointer;
    transition: box-shadow 0.3s;
}

.mermaid-svg-wrapper {
    border-radius: 10px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    background-color: white;
}

.mermaid-svg-wrapper svg {
    border-radius: 8px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.mermaid-container:hover {
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
}

.zoom-wrapper {
    width: 100%;
    height: 80vh;
    overflow: hidden;
    position: relative;
    border: 1px solid #ddd;
    background: #fff;
}

.zoom-inner {
    transform-origin: 0 0;
    cursor: grab;
}
</style>
