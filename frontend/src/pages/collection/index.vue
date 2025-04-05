<template>
    <div>
        <div>
            <a-flex gap="middle" horizontal>
                <a-button @click="add">
                    <PlusOutlined />新增
                </a-button>
                <a-input placeholder="集合名称" v-model:value="query.collectionName"></a-input>
                <a-button @click="queryData">
                    <SearchOutlined />查询
                </a-button>
            </a-flex>
        </div>
        <a-divider style="margin-top: 10px;margin-bottom: 10px;" />
        <div>
            <a-list item-layout="horizontal" :data-source="records">
                <template #renderItem="{ item }">
                    <a-list-item style="width: 100%;">
                        <template #actions>
                            <a key="list-loadmore-more" @click="fileChange(item)">上传文档</a>
                            <a key="list-loadmore-more" @click="del(item)">删除</a>
                        </template>
                        <a-list-item-meta>
                            <template #title>
                                <strong>{{ item }}</strong>
                            </template>
                        </a-list-item-meta>
                    </a-list-item>
                </template>
            </a-list>
        </div>
        <input v-show="false" multiple :disabled="isLoading" ref="inputFile" type="file" @change="addFile($event)" />
    </div>
</template>
<script>
import { message, Modal } from 'ant-design-vue';
import { Delete, GetList, Upload } from '../../../wailsjs/go/biz/CollectionBiz.js';


export default {
    mounted() {
        this.query = this.$store.queryData["collection"] || {};
        this.getList();
    },
    data() {
        return {
            query: {},
            records: [],
            record: {},
            file: null,

        };
    },
    methods: {
        fileChange(record) {
            this.record = record;
            this.$refs.inputFile.click();
        },
        addFile(event) {
            this.file = event.target.files;
            this.convertFilesToBase64(event);
        },
        // 将文件转换为 Base64 字符串
        async convertFilesToBase64(event) {
            const files = this.file;
            const promises = {};
            // 遍历每个文件，将其转换为 Base64
            for (let i = 0; i < files.length; i++) {
                let f = files[i]
                promises[f.name] = await this.fileToBase64(f)
            }
            this.uploadFile(promises);
        },
        async fileToBase64(file) {
            return new Promise((resolve, reject) => {
                const reader = new FileReader();
                reader.onloadend = () => {
                    resolve(reader.result); // 读取文件为 Base64 格式
                };
                reader.onerror = reject;
                reader.readAsDataURL(file); // 将文件读取为 Data URL（Base64 格式）
            });
        },
        uploadFile(base64Files) {
            Upload(this.record, base64Files).then(res => {
                if (res.code == 200) {
                    message.info("上传成功");
                } else {
                    message.error(res.msg);
                }
            }).catch(err => {
                message.error(err);
            });
        },
        pageClick(pagination) {
            this.pages = pagination;
            this.getList();
        },
        customRow(record) {
            return {
                onClick: (event) => {

                },
                onDblclick: (event) => {
                    this.record = record;
                },
                onContextmenu: (event) => { },
                onMouseenter: (event) => { },  // 鼠标移入行
                onMouseleave: (event) => { }
            };
        },
        getList() {
            GetList(this.query).then(res => {
                this.records = res.data || [];
            })
        },
        queryData() {
            this.$store.setQueryData("collection", this.query);
            this.records = [];
            this.getList();
        },
        add(data) {
            this.$store.setAction("add");
            this.$router.push({
                path: "/collection/form",
                query: {
                    id: data.id,
                    system: this.activeKey,
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
                    Delete(data).then(res => {
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