<template>
    <div>
        <div ref="chartContainer" class="echarts-chart"></div>
        <a-flex align="center" justify="center"><a-button type="link" @click="changeChart" :disabled="$store.loading">更换图表</a-button></a-flex>
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
    data() {
        return {
            isDisabled: false,
        }
    },
    methods: {
        changeChart() {
            this.$store.setLoading(true);
            let message = "这是一个echarts图表的JSON配置数据,请给我更新一种展现方式,并使用md代码块输出，语言类型为：chart。\n不要输出额外的信息。 图表json数据：" + JSON.stringify(this.chartData);
            this.$bus.emit("changeChart", message);
            setTimeout(() => {
                this.$store.setLoading(false);
            }, 3000);
        },
        renderChart() {
            if (!this.$refs.chartContainer || !this.chartData) return;

            if (this.chartInstance) {
                this.chartInstance.dispose();
                this.chartInstance = null;
            }

            this.chartInstance = echarts.init(this.$refs.chartContainer);
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
