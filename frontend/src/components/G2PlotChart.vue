<template>
    <div>
        <div ref="chartContainer" class="g2plot-chart"></div>
    </div>
</template>

<script>
import * as G2Plot from "@antv/g2plot";

export default {
    props: {
        chartData: {
            type: Object,
            required: true
        }
    },
    mounted() {
        this.$nextTick(() => {
            this.renderChart();
        });
    },
    methods: {
        renderChart() {
            if (!this.$refs.chartContainer) return;

            // **销毁旧实例，避免重叠**
            if (this.chartInstance) {
                this.chartInstance.destroy();
                this.chartInstance = null;
            }

            // **获取正确的 G2Plot 图表类**
            const { type = "Column", ...config } = this.chartData;
            const ChartClass = G2Plot[type];

            if (!ChartClass) {
                console.warn(`G2Plot 不支持的图表类型: "${type}"，已回退到 "Column"`);
                return;
            }

            // **创建并渲染图表**
            this.chartInstance = new ChartClass(this.$refs.chartContainer, config);
            this.chartInstance.render();
        }
    },
    watch: {
        chartData: {
            deep: true,
            handler() {
                this.renderChart();
            }
        }
    },
    beforeUnmount() {
        if (this.chartInstance) {
            this.chartInstance.destroy();
            this.chartInstance = null;
        }
    }
};
</script>

<style>
.g2plot-chart {
    width: 100%;
    height: 350px;
    padding: 10px;
}

</style>
