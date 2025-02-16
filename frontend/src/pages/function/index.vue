<template>
    <div>
        <div>
            <a-flex gap="middle" horizontal>
                <a-button @click="add">
                    <PlusOutlined />新增
                </a-button>
                <a-input placeholder="函数名称" v-model:value="query.funcName"></a-input>
                <a-button @click="queryData">
                    <SearchOutlined />查询
                </a-button>
            </a-flex>
        </div>
        <a-divider style="margin-top: 10px;margin-bottom: 10px;" />
        <div>
            <a-table ref="table" :data-source="records" :columns="columns" @change="pageClick" :pagination="pages"
                row-key="id" :scroll="{ x: 300,  y: 500 }" :customRow="customRow">
                <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'action'">
                        <a-space warp>
                            <a key="list-loadmore-more" @click="copy(record)">复制</a>
                            <a key="list-loadmore-more" @click="detail(record)">查看</a>
                            <a key="list-loadmore-more" @click="edit(record)">编辑</a>
                            <a key="list-loadmore-more" @click="del(record)">删除</a>
                        </a-space>
                    </template>
                </template>
            </a-table>
        </div>
    </div>
</template>
<script>
import { message, Modal } from 'ant-design-vue';
import { Delete, Page } from '../../../wailsjs/go/biz/FunctionBiz.js';


export default {
    mounted() {
        this.query = this.$store.queryData["function"] || {};
        this.loadPage();
    },
    data() {
        return {
            query: {},
            pages: {
                current: 1,
                size: 10,
                total: 0,
                showQuickJumper: true,
                pageSizeOptions: ["10", "20", "50", "100", "200"]
            },
            records: [],
            record: {},
            columns: [
                { title: "函数名称", dataIndex: "funcName", width: 300 },
                { title: "函数内容", dataIndex: "funcContent", ellipsis: true },
                { title: '操作', key: 'action', fixed: 'right', width: 180 }
            ],
        };
    },
    methods: {
        pageClick(pagination) {
            this.page = pagination;
            this.loadPage();
        },
        customRow(record) {
            return {
                onClick: (event) => {

                },
                onDblclick: (event) => {
                    message.success("选中行新增")
                    this.record = record;
                },
                onContextmenu: (event) => { },
                onMouseenter: (event) => { },  // 鼠标移入行
                onMouseleave: (event) => { }
            };
        },
        loadPage() {
            Page({ page: this.pages, ...this.query }).then(res => {
                let data = res.data || {};
                if (res.code == 200) {
                    this.records = data.rows;
                } else {
                    this.$message.error(res.msg)
                }
            })
        },
        queryData() {
            this.$store.setQueryData("function", this.query);
            this.pages.current = 1;
            this.records = [];
            this.loadPage();
        },
        add(data) {
            this.$store.setAction("add");
            this.$router.push({
                path: "/function/form",
                query: {
                    id: data.id,
                    system: this.activeKey,
                },
            });
        },
        copy(data) {
            this.$store.setAction("add");
            this.$router.push({
                path: "/function/form",
                query: {
                    id: data.id,
                },
            });
        },
        edit(data) {
            this.$store.setAction("edit");
            this.$router.push({
                path: "/function/form",
                query: {
                    id: data.id,
                },
            });
        },
        detail(data) {
            this.$store.setAction("detail");
            this.$router.push({
                path: "/function/form",
                query: {
                    id: data.id,
                },
            });
        },
        del(data) {
            let m = this;
            Modal.confirm({
                title: '确定删除？',
                okText: '确定',
                cancelText: '取消',
                onOk() {
                    Delete(data.id).then(res => {
                        if (res.code == 200) {
                            message.success('删除成功!');
                            m.queryData();
                        } else {
                            message.error(res.msg);
                        }
                    });
                },
                onCancel() {
                },
            });
        },
    },
};
</script>