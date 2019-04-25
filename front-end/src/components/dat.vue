<template>
  <section class="bg-white">


    <el-row :gutter="20">

      <el-col :span="8">
        <div class="intro">
          <img src="../assets/images/music.png">
          <div>
            <br>
            <b> Create and sell your dat and avatar</b>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="wallet-container">
          <div class=" bridge-item-container">
            <div class="shop-avatar-info-container">
              <div class="role-info-header">
                Your current number of dat
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
                show-file-list
                :action="upLoadDatPath"
                :data="upLoadAdditionalData"
                :file-list="uploadFileList"
                :auto-upload="false">
                <el-button slot="trigger" size="small" type="primary">select file</el-button>
                <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">upload to
                  server
                </el-button>
                <div class="el-upload__tip" slot="tip">jpg/png files with a size less than 500kb</div>
              </el-upload>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="upload-file-container">
          <el-table
            :data="tableData"
            stripe
            style="width: 100%">
            <el-table-column
              prop="fileName"
              label="File Name"
              width="180">
            </el-table-column>
            <el-table-column
              prop="uploadTime"
              label="Upload Time"
              width="180">
            </el-table-column>
            <el-table-column
              prop="fileSize"
              label="File Size">
            </el-table-column>
          </el-table>

        </div>
      </el-col>

    </el-row>

  </section>
</template>
<script>

  export default {
    props: {},
    data() {
      return {
        fileList: [],
        upLoadDatPath: undefined,
        upLoadAdditionalData: undefined,
        httpPath: undefined,
        wallet: {
          publicBalance: undefined,
          balance: undefined,
          publicEther: undefined
        },
        // uploadFileList: [],
        // A fake data list
        tableData: [
          {
            fileName: 'a.jpg',
            uploadTime: '2017-05-03 17:30',
            fileSize: '3 MB'
          }, {
            fileName: 'df.jpg',
            uploadTime: '2007-05-03 17:30',
            fileSize: '4 MB'
          }, {
            fileName: 'dfs.mp3',
            uploadTime: '2017-12-03 17:30',
            fileSize: '45 MB'
          }, {
            fileName: 'this.jpg',
            uploadTime: '2017-05-03 04:24',
            fileSize: '567 KB'
          }, {
            fileName: 'c.png',
            uploadTime: '2017-05-03 17:30',
            fileSize: '3 MB'
          }]
      }
    },
    methods: {
      submitUpload: function () {
        // let fileNameList = this.$refs.upload.fileList
        // console.log(this.uploadFileList[0])
        // let currentTime = new Date()
        // let uploadTime = currentTime.getFullYear() + '-' + currentTime.getMonth() + '-' + currentTime.getDate() + ' '
        //   + currentTime.getHours() + ':' + currentTime.getMinutes()
        // fileNameList.forEach(function (fileName) {
        //   this.tableData.push({
        //     fileName: fileName,
        //     uploadTime: uploadTime,
        //     fileSize:
        //   })
        // })
        this.$refs.upload.submit();
      },
    },
    beforeMount: function () {
      // console.log(this.$store.state.account)
      this.upLoadAdditionalData = {
        address: this.$cookies.get('account').address
      }
      this.httpPath = this.$store.state.config.httpPath
      this.upLoadDatPath = this.httpPath + "/file/dat"
    }
  }
</script>
