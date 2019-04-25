<template>
    <section class="bg-light">
      <div class="row-container">
        <div class="intro text-center">
          <img src="../assets/images/avatar.jpg">
          <div>
            <!--<h3> Market </h3>-->
            <!--<br>-->
            <!--<br>-->
            <!--<b> With bridge</b>-->
            <br>
            <b> Create and sell your dat and avatar</b>
            <!--<br>-->
            <!--<b>between public and private chain</b>-->
          </div>
        </div>

        <div class="wallet-container">
          <div class=" bridge-item-container">
            <div class="shop-avatar-info-container">
              <div class="role-info-header">
                Your current number of avatar
              </div>
              <div>
                1
              </div>
            </div>
          </div>
          <div class="bridge-item-container">
            <div class="token-container">
              <el-upload
                class="upload-demo upload-dat"
                ref="upload"
                name="file"
                :action="upLoadDatPath"
                :data="upLoadAdditionalData"
                :auto-upload="false">
                <el-button slot="trigger" size="small" type="primary">select file</el-button>
                <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">upload to server</el-button>
                <div class="el-upload__tip" slot="tip">jpg/png files with a size less than 500kb</div>
              </el-upload>
            </div>
          </div>
        </div>
      </div>
    </section>
</template>
<script>
  export default {
    props: {
    },
    data () {
      return {
        fileList: [],
        upLoadDatPath: undefined,
        upLoadAdditionalData:undefined,
        httpPath:undefined,
        wallet: {
          publicBalance: undefined,
          balance: undefined,
          publicEther: undefined
        },
      }
    },
    methods: {
      submitUpload: function () {
        this.$refs.upload.submit();
      },
    },
    beforeMount: function() {
      // console.log(this.$store.state.account)
      this.upLoadAdditionalData = {
        address: this.$cookies.get('account').address
      }
      this.httpPath = this.$store.state.config.httpPath
      this.upLoadDatPath = this.httpPath + "/file/avatar"
    }
  }
</script>
