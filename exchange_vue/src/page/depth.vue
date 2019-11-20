<template>
    <div class="fillcontain">
        <head-top></head-top>
        <div class="table_container">
            <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
                <el-tab-pane label="OKEX" name="okex"></el-tab-pane>
                <el-tab-pane label="BITMEX" name="bitmex"></el-tab-pane>
                <el-tab-pane label="火币" name="huobi"></el-tab-pane>
                <el-tab-pane label="币安" name="binance"></el-tab-pane>

                <div style="float:left;margin-left: 10px;margin-bottom: 10px">
                    <el-button type="primary" @click="handleEdit()" >深度编辑</el-button>
                </div>

                <div style="float:right;margin-right: 10px;margin-bottom: 10px">
                    <el-button type="warning" @click="handleRestart('local')" >重启本机</el-button>
                    <el-button type="danger" @click="handleRestart('other')" >重启其他</el-button>
                </div>

                <el-dialog title="编辑平台数据" :visible.sync="dialogFormVisible">
                    <el-form :model="depth" :data="tableData">
                        <el-form-item label="币种对" label-width="100px" v-for="(i,index) in tableData">
                            <el-col :span="11"style="margin-left: 10px">
                                <el-input placeholder="输入币种对" v-model="i.symbol"></el-input>
                            </el-col>
                            <el-col :span="1" style="margin-left: 10px">
                                <i @click="addRow()" class="el-icon-plus"></i>
                            </el-col>
                            <el-col :span="1">
                                <i @click="reduceRow(index)" class="el-icon-minus"></i>
                            </el-col>
                        </el-form-item>
                    </el-form>
                    <div slot="footer" class="dialog-footer">
                        <el-button @click="dialogFormVisible = false">取 消</el-button>
                        <el-button type="primary" @click="WsConnect()">WS连接测试</el-button>
                    </div>
                </el-dialog>

                <el-dialog title="ws连接检测" :visible.sync="dialogTableVisible">
                    <el-table :data="wsdata">
                        <el-table-column label="序号" width="100">
                            <template slot-scope="scope">
                                {{scope.$index+1}}
                            </template>
                        </el-table-column>
                        <el-table-column property="symbol" label="币种对"></el-table-column>
                        <el-table-column property="status" label="检测状态">
                            <template slot-scope="scope">
                                <el-row>
                                    <i class="el-icon-loading" v-if="scope.row.status == 0"></i>
                                    <i class="el-icon-check" v-if="scope.row.status == 1"></i>
                                </el-row>
                            </template>
                        </el-table-column>
                    </el-table>
                    <div slot="footer" class="dialog-footer">
                        <el-button @click="dialogTableVisible = false">取 消</el-button>
                        <el-button type="primary" @click="commit()">确认</el-button>
                    </div>
                </el-dialog>

            <el-table ref="filterTable" :data="tableData" style="width: 100%" >
                <el-table-column label="序号" width="100">
                    <template slot-scope="scope">
                        {{scope.$index+1}}
                    </template>
                </el-table-column>
                <el-table-column prop="symbol" label="币种对" sortable width="180"></el-table-column>
                <el-table-column prop="price" label="价格" sortable width="180"></el-table-column>
                <el-table-column
                    prop="tag"
                    label="类型"
                    width="180"
                    :filters="[{ text: '现货', value: 0 }, { text: '期货', value: 1 }]"
                    :filter-method="filterTag"
                    filter-placement="bottom-end">
                    <template slot-scope="scope">
                        <el-tag type="primary" v-if="scope.row.tag == 0">现货</el-tag>
                        <el-tag type="success" v-if="scope.row.tag == 1">期货</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="time"  label="时间"></el-table-column>
                <el-table-column prop="status" label="状态">
                    <template slot-scope="scope">
                        <el-row>
                            <el-tag type="danger" v-if="scope.row.status == 0">已停止</el-tag>
                            <el-tag type="success" v-if="scope.row.status == 1">运行中</el-tag>
                        </el-row>
                    </template>
                </el-table-column>
            </el-table>
            </el-tabs>
        </div>
    </div>
</template>

<script>
    import headTop from '../components/headTop'
    export default {
        data() {
            return {
                depth:{},
                activeName:"okex",
                dialogFormVisible:false,
                dialogTableVisible:false,
                indexId: '',
                checkId: '',
                wsdata:[],
                tableData: []
            }
        },
        components: {
            headTop,
        },
        mounted:function() {
            this.indexId = setInterval(this.getList,1000)
        },
        methods: {
            filterTag(value, row) {
                return row.tag === value;
            },
            // 编辑
            handleEdit() {
                clearInterval(this.indexId);
                this.dialogFormVisible = true;
            },

            WsConnect() {
                this.dialogFormVisible = false;
                this.dialogTableVisible = true;
                this.sendWsConnect();
                this.checkId = setInterval(this.WsCheck,1000)
            },

            // 点击切换平台
            handleClick(tab, event) {
                this.$options.methods.getList.bind(this)()
            },

            // 增加一行
            addRow() {
                this.tableData.push({
                    symbol:""
                })
            },
            // 删除一行
            reduceRow(index) {
                this.tableData.splice(index,1)
            },

            // list 请求
            getList() {
                var url = '/depth?platform=' + this.activeName;
                this.$http.get(url)
                    .then(response => {
                        if (response.data.status === 1) {
                            this.tableData = response.data.data
                        }
                    })
            },
            WsCheck() {
                var url = '/depth/check/' + this.activeName;
                this.$http.get(url)
                    .then(response => {
                        this.wsdata = response.data.data
                    })
            },

            sendWsConnect() {
                var url = '/depth/' + this.activeName;

                var data = [];
                for (var i = 0;i< this.tableData.length;i++) {
                    data[i] = this.tableData[i].symbol
                }
                this.$http.put(url,{symbols:data})
                    .then(response => {
                        console.log(response)
                    })
            },

            commit() {
                clearInterval(this.checkId);
                var url = '/depth/commit?platform=' + this.activeName;
                this.$http.get(url)
                    .then(response => {
                        if (response.data.status == true) {
                            this.$message.success('配置提交成功,在重启ws客户端进程后生效');
                            this.dialogTableVisible = false;
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
