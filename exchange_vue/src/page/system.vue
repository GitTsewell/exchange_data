<template>
    <div class="fillcontain">
        <head-top></head-top>
        <div class="table_container">
            <el-table :data="tabledatas" border>
                <el-table-column label="序号" width="100">
                    <template slot-scope="scope">
                        {{scope.$index+1}}
                    </template>
                </el-table-column>
                <el-table-column label="名称">
                    <template slot-scope="scope">
                        <el-input placeholder="请输入内容" v-show="scope.row.show" v-model="scope.row.name"></el-input>
                        <span v-show="!scope.row.show">{{scope.row.name}}</span>
                    </template>
                </el-table-column>
                <el-table-column label="指令">
                    <template slot-scope="scope">
                        <el-input placeholder="请输入内容" v-show="scope.row.show" v-model="scope.row.action"></el-input>
                        <span v-show="!scope.row.show">{{scope.row.action}}</span>
                    </template>
                </el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button @click="scope.row.show =true">编辑</el-button>
                        <el-button @click="commit(scope.row)">保存</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>

<script>
    import headTop from '../components/headTop'
    export default {
        components: {
            headTop,
        },

        data() {
            return {
                tabledatas: [],
            }
        },

        mounted:function() {
            this.$http.get('/system')
                .then(response => {
                    this.tabledatas = response.data.data
                })
        },

        methods:{
            commit(index) {
                index.show = false;
                this.$http.put("/system",{key:index.key,action:index.action})
                    .then(response => {
                        if (response.data.status == 1) {
                            this.$message.success("编辑成功")
                        }
                    })
            }
        }
    }
</script>

<style lang="less">
    @import '../style/mixin';
    .table_container{
        padding: 20px;
    }
</style>
