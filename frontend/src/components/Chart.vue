<template>
    <div>
        <div ref="chartContainer" class="echarts-chart"></div>
    </div>
</template>

<script>
import * as echarts from 'echarts';

export default {
    props: {
        chartData: {
            type: Object,
            required: true
        }
    },
    mounted() {
        this.$nextTick(() => {
            this.renderChart();  // 初始化时渲染图表
        });
    },
    methods: {
        renderChart() {
            if (!this.$refs.chartContainer || !this.chartData) return;

            // **销毁旧实例，避免重叠**
            if (this.chartInstance) {
                this.chartInstance.dispose();
                this.chartInstance = null;
            }

            // **初始化 ECharts 实例**
            this.chartInstance = echarts.init(this.$refs.chartContainer);

            // **使用后台传递的完整配置进行渲染**
            this.chartInstance.setOption(this.chartData);
        }
    },
    watch: {
        chartData: {
            deep: true,
            handler() {
                this.renderChart();  // 如果数据变化，重新渲染图表
            }
        }
    },
    beforeUnmount() {
        if (this.chartInstance) {
            this.chartInstance.dispose();
            this.chartInstance = null;
        }
    }
};
</script>

<style>
.echarts-chart {
    width: 100%;
    height: 350px;
    padding: 10px;
}
</style>
