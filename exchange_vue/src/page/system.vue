<template>
    <div class="fillcontain">
        <head-top></head-top>
        <div class="table_container">
            <el-table
                :data="tableData"
                border
                style="width: 100%">
                <el-table-column label="序号" width="100">
                    <template slot-scope="scope">
                        {{scope.$index+1}}
                    </template>
                </el-table-column>
                <el-table-column
                    prop="instruct"
                    label="指令名称"
                    width="180">
                </el-table-column>
                <el-table-column
                    prop="exec"
                    label="指令内容">
                </el-table-column>
                <el-table-column
                    label='hah'
                    class-name="column-bg-color-editable"
                    width="100"
                    show-overflow-tooltip>
                    <template scope="scope">
                        <div class="input-box">
                            <el-input size="small" @blur="handleInputBlur" @cell-click ="handleCellClick" v-model="scope.row.description" ></el-input>
                        </div>
                        <span>{{scope.row.description}}</span>
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
                tableData: [{
                    instruct: '重启本机ws客户端',
                    exec: 'supervisorctl restart a1',
                    description: "dfdf"
                }, {
                    instruct: '重启其他ws客户端',
                    exec: 'ssh xx.xx.xx.x supervisorctl restart a1',
                    description: "dsdsdsdfdf"
                }]
            }
        },

        methods: {
            //单元格点击后，显示input，并让input 获取焦点
            handleCellClick:function(row, column, cell, event){
                emptransfer.addClass(cell,'current-cell');
                if(emptransfer.getChildElement(cell,3) !== 0){
                    var _inputParentNode =emptransfer.getChildElement(cell,3);
                    if(_inputParentNode.hasChildNodes()&& _inputParentNode.childNodes.length > 2) {
                        var _inputNode = _inputParentNode.childNodes[2];
                        if(_inputNode.tagName === 'INPUT'){
                            _inputNode.focus();
                        }
                    }
                }
            },
            //input框失去焦点事件
            handleInputBlur:function(event){   //当 input 失去焦点 时,input 切换为 span，并且让下方 表格消失（注意，与点击表格事件的执行顺序）
                var _event = event;
                setTimeout(function(){
                    var _inputNode = _event.target;
                    if(emptransfer.getParentElement(_inputNode,4)!==0){
                        var _cellNode = emptransfer.getParentElement(_inputNode,4);
                        emptransfer.removeClass(_cellNode,'current-cell');
                        emptransfer.removeClass(_cellNode,'current-cell2');
                    }
                },200);
            }
        }
    }
</script>

<style lang="less">
    @import '../style/mixin';
    .table_container{
        padding: 20px;
    }

    .tb-edit .input-box {
        display: none
    }
    .tb-edit .current-cell .input-box {
        display: block;
        margin-left: -15px;
    }
</style>
