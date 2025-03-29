<template>
    <div v-html="renderedMessage"></div>
</template>

<script>
export default {
    props: {
        renderedMessage: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            // 可以根据需要添加一些状态，如loading标识等
        };
    },
    methods: {
        // 如果需要，可以在这里处理一些额外的逻辑
    },
};
</script>