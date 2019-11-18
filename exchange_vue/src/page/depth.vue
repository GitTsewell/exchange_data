<template>
    <div class="fillcontain">
        <head-top></head-top>
        <div class="table_container">
            <el-table ref="filterTable" :data="tableData" style="width: 100%">
                <el-table-column prop="id" label="#" width="180"></el-table-column>
                <el-table-column prop="symbol" label="币种对" width="180"></el-table-column>
                <el-table-column prop="price" label="价格" width="180"></el-table-column>
                <el-table-column
                    prop="type"
                    label="类型"
                    width="100"
                    :filters="[{ text: '现货', value: '现货' }, { text: '期货', value: '期货' }]"
                    :filter-method="filterTag"
                    filter-placement="bottom-end">
                    <template slot-scope="scope">
                        <el-tag
                            :type="scope.row.type === '现货' ? 'primary' : 'success'"
                            disable-transitions>{{scope.row.type}}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="time" label="时间"></el-table-column>
                <el-table-column prop="status" label="状态">
                    <template slot-scope="scope">
                        <el-row>
                            <el-tag type="danger" v-if="scope.row.status == 0">已停止</el-tag>
                            <el-tag type="success" v-if="scope.row.status == 1">运行中</el-tag>
                        </el-row>
                    </template>
                    <el-tag type="success">标签二</el-tag>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>

<script>
    import headTop from '../components/headTop'
    export default {
        data() {
            return {
                tableData: [{
                    id:1,
                    symbol: 'BCH-USDT',
                    price: '350',
                    type: '现货',
                    time: '2019-11-18T17:25:29.874+08:00',
                    status: 1
                }, {
                    id:2,
                    symbol: 'BCH-USDT',
                    price: '350',
                    type: '现货',
                    time: '2019-11-18T17:25:29.874+08:00',
                    status: 0
                },{
                    id:3,
                    symbol: 'BCH-USD-191227',
                    price: '350',
                    type: '期货',
                    time: '2019-11-18T17:25:29.874+08:00',
                    status: 1
                },{
                    id:4,
                    symbol: 'BCH-USDT',
                    price: '350',
                    type: '现货',
                    time: '2019-11-18T17:25:29.874+08:00',
                    status: 0
                }]
            }
        },
        components: {
            headTop,
        },
        methods: {
            resetDateFilter() {
                this.$refs.filterTable.clearFilter('date');
            },
            clearFilter() {
                this.$refs.filterTable.clearFilter();
            },
            formatter(row, column) {
                return row.address;
            },
            filterTag(value, row) {
                return row.type === value;
            },
            filterHandler(value, row, column) {
                const property = column['property'];
                return row[property] === value;
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
