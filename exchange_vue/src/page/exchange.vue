<template>
    <div class="fillcontain">
        <head-top></head-top>
        <div class="table_container">
            <el-form :model="ruleForm" label-width="80px" inline>
                <el-form-item label="OKEX" prop="okex" required>
                    <el-switch v-model="ruleForm.okex"></el-switch>
                </el-form-item>
                <el-form-item label="BITMEX" prop="bitmex" required>
                    <el-switch v-model="ruleForm.bitmex"></el-switch>
                </el-form-item>
                <el-form-item label="火币" prop="huobi" required>
                    <el-switch v-model="ruleForm.huobi"></el-switch>
                </el-form-item>
                <el-form-item label="币安" prop="binance" required>
                    <el-switch v-model="ruleForm.binance"></el-switch>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="submitForm()">更新</el-button>
                    <el-button @click="resetForm('ruleForm')">重置</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
    import headTop from '../components/headTop'

    export default {
        data() {
            return {
                ruleForm: {
                    okex:true,
                    bitmex:true,
                    huobi:true,
                    binance:true,
                },
            };
        },

        components: {
            headTop,
        },

        mounted:function() {
            this.edit()
        },

        methods: {
            edit() {
                this.$http.get("/exchange")
                    .then(response => {
                        this._data.ruleForm = response.data.data
                    })
            },

            submitForm() {
                this.$http.put("/exchange",this._data.ruleForm)
                    .then(response => {
                        if (response.data.status == 1) {
                            this.$message.success("更新成功")
                        }
                    })
            },
            resetForm(formName) {
                this.$refs[formName].resetFields();
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


