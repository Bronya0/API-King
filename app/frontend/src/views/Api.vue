<template>
  <n-flex vertical justify="start" style="text-align: left;">
    <!-- 接口调试区域 -->
    <n-card>
      <n-flex align="center">
        <n-select
            v-model:value="form.groupName"
            :options="groupOptions"
            placeholder="分组"
            style="width: 10%; "
            clearable
            filterable
            tag
        />
        <n-input v-model:value="form.apiName" placeholder="名称" style="width: 10%; "/>
        <n-select v-model:value="form.method" :options="methodOptions" placeholder="选择方法" style="width: 100px; "/>
        <n-input v-model:value="form.url" placeholder="URL" style="width: 56%; "/>
        <n-button :disabled="!form.url" :loading="loading.request" @click="sendRequest" style="width: 88px">
          发送请求
        </n-button>
        <n-button @click="clear">清理</n-button>

      </n-flex>
    </n-card>
    <n-card>
      <n-tabs v-model:value="activeTab" type="card" style="text-align: left; ">

        <n-tab-pane display-directive="show" name="请求头" tab="请求头">
          <AceEditor
              ref="header_editor"
              mode="json"
              style="height: 400px; "
          />
        </n-tab-pane>

        <n-tab-pane display-directive="show" name="查询参数" tab="查询参数">
          <AceEditor
              ref="query_editor"
              mode="json"
              style="height: 400px; "
          />
        </n-tab-pane>

        <n-tab-pane display-directive="show" name="请求体" tab="请求体">
          <n-flex vertical>
            <n-flex align="center">
              <n-select v-model:value="bodyType" @update:value="body_editor.setValue(null)" :options="bodyTypeOptions" style="width: 120px; "
                        placeholder="选择格式" />
              form下，body的json里的key以@开头，值为文件路径，则表示文件上传（注意路径用\\ or /）
            </n-flex>
            <AceEditor
                ref="body_editor"
                mode="json"
                style="height: 400px; "
            />
          </n-flex>
        </n-tab-pane>

        <n-tab-pane display-directive="show" name="响应" tab="响应">
          <n-flex vertical>
            耗时：{{ responseTime }}ms | status：{{ responseStatus }} | size：{{ responseSize }} bytes
            <AceEditor
                ref="resp_editor"
                mode="json"
                style="height: 600px; "
            />
          </n-flex>
        </n-tab-pane>

        <n-tab-pane display-directive="show" name="历史记录" tab="历史记录">
          <n-flex vertical>
            <!-- 添加搜索区域 -->
            <n-flex align="center" style="margin-bottom: 16px;">
              <n-select
                  v-model:value="searchForm.groupName"
                  :options="groupOptions"
                  placeholder="选择分组"
                  style="width: 150px; margin-right: 8px;"
                  clearable
              />
              <n-input
                  v-model:value="searchForm.apiName"
                  placeholder="接口名称"
                  style="width: 200px; margin-right: 8px;"
                  clearable
              />
              <n-input
                  v-model:value="searchForm.url"
                  placeholder="URL"
                  style="width: 300px; margin-right: 8px;"
                  clearable
              />
              <n-button @click="fetchHistory">查询</n-button>
            </n-flex>

            <n-data-table
                ref="historyTableRef"
                :columns="refColumns(historyColumns)"
                :data="historyList"
                :loading="loading.history"
                :pagination="historyPagination"
                size="small"
            />
          </n-flex>
        </n-tab-pane>

        <n-tab-pane name="压测" tab="压测">
          <n-flex align="center">
            并发量
            <n-input-number
                v-model:value="benchmarkForm.clientNum"
                style="width: 150px;"
            />
            压测时间
            <n-input-number
                v-model:value="benchmarkForm.secs"
                style="width: 100px;"
            />
            超时
            <n-input-number
                v-model:value="benchmarkForm.timeout"
                style="width: 100px;"
            />
            <n-button :disabled="!form.url" @click="benchmark" :loading="loading.benchmark">开始压测</n-button>
          </n-flex>
          <n-code :code="benchmarkForm.resp" language="json" show-line-numbers/>
        </n-tab-pane>

      </n-tabs>
    </n-card>


  </n-flex>

</template>

<script setup>
import {h, onMounted, ref} from 'vue';
import {NButton, NDataTable, NFlex, NInput, NSelect, useMessage} from 'naive-ui';
import {Benchmark, ProxyInsert, ProxyQuery, ProxySql, ProxyWithInfo} from '../../bindings/app/backend/service/Api';
import timeutil from "../utils/timeutil";
import AceEditor from "../components/AceEditor.vue";
import {refColumns} from "../utils/common";

const message = useMessage();

// 数据
const form = ref({
  groupName: null,
  apiName: '',
  method: 'GET',
  url: '',
});

const benchmarkForm = ref({
  clientNum: 10,
  secs: 5,
  timeout: 5,
  resp: "",
});

const header_editor = ref(null);
const query_editor = ref(null);
const body_editor = ref(null);
const resp_editor = ref(null);

const activeTab = ref('响应');


// 请求体类型
const bodyType = ref('json');
const bodyTypeOptions = [
  {label: 'JSON', value: 'json'},
  {label: 'Form', value: 'form'},
];

// 搜索表单
const searchForm = ref({
  groupName: null,
  apiName: '',
  url: ''
});

// 分组下拉选项
const groupOptions = ref([]);

// 获取分组选项的函数
const fetchGroupOptions = async () => {
  try {
    const sql = `
      SELECT DISTINCT group_name
      FROM api
      WHERE group_name IS NOT NULL AND group_name != ''
      ORDER BY group_name
    `;
    const groups = await ProxyQuery(sql);
    groupOptions.value = groups.map(group => ({
      label: group.group_name,
      value: group.group_name
    }));
  } catch (e) {
    message.error('获取分组列表失败: ' + e.message);
    console.error(e.message)

  }
};

// 请求方法选项
const methodOptions = [
  {label: 'GET', value: 'GET'},
  {label: 'POST', value: 'POST'},
  {label: 'PUT', value: 'PUT'},
  {label: 'HEAD', value: 'HEAD'},
  {label: 'PATCH', value: 'PATCH'},
  {label: 'OPTIONS', value: 'OPTIONS'},
  {label: 'DELETE', value: 'DELETE'},
];

// 响应数据
const responseTime = ref(0);
const responseStatus = ref(0);
const responseSize = ref(0);
const loading = ref({
  request: false,
  history: false,
  benchmark: false,
});

// 历史记录分页
const historyPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page) => {
    historyPagination.value.page = page;
    fetchHistory();
  },
  onUpdatePageSize: (pageSize) => {
    historyPagination.value.pageSize = pageSize;
    historyPagination.value.page = 1;
    fetchHistory();
  },
  itemCount: 0,
  suffix: ({itemCount}) => `共 ${itemCount} 条记录`,
});

// 历史记录列表
const historyList = ref([]);
const historyTableRef = ref();

// 历史记录列定义
const historyColumns = [
  {title: '分组', key: 'group_name', width: 100},
  {title: '接口名称', key: 'api_name', width: 200},
  {title: '方法', key: 'method', width: 80},
  {title: 'URL', key: 'url', width: 300},
  {title: '请求头', key: 'headers', width: 100},
  {title: '查询参数', key: 'params', width: 200},
  {title: 'Type', key: 'type', width: 80},
  {title: '请求体', key: 'body', width: 100},
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row) =>
        h(NFlex, null, {
          default: () => [
            h(NButton, {
              size: 'small',
              onClick: () => fillForm(row),
            }, {default: () => '使用'}),
            h(NButton, {
              size: 'small',
              type: 'error',
              secondary: true,
              onClick: () => deleteHistory(row),
            }, {default: () => '删除'}),
          ]
        })
  },
  {
    title: '使用时间',
    key: 'updated_at',
    width: 150,
    render: (row) => {
      return row.updated_at ? timeutil.formattedStringToLocalStr(row.updated_at) : '';
    },
  }
];


// 发送请求
const sendRequest = async () => {
  loading.value.request = true;
  try {
    activeTab.value = '响应';
    resp_editor.value.setValue('请求中……');

    let {
      headersObj,
      queryObj,
      requestParams,
    } = await prepareSendRequest()

    const startTime = Date.now();

    let resp = await ProxyWithInfo(
        form.value.method,
        form.value.url,
        headersObj,
        queryObj,
        requestParams.json || "",
        requestParams.form || {},
        requestParams.files || {},
        false // useToken
    );
    responseTime.value = Date.now() - startTime;

    responseStatus.value = resp['status'];
    responseSize.value = resp['size'];
    let body_str = resp['body'];
    // json格式化
    if (body_str && body_str.startsWith('{')) {
      body_str = JSON.stringify(JSON.parse(body_str), null, 2);
    }
    resp_editor.value.setValue(body_str || "");

    // 保存到历史记录
    let instances = [{
      group_name: form.value.groupName,
      api_name: form.value.apiName,
      method: form.value.method,
      url: form.value.url,
      headers: JSON.stringify(headersObj),
      params: JSON.stringify(queryObj),
      body: requestParams.json || requestParams.form,
      type: bodyType.value,
    }];
    await ProxyInsert("api", instances, ['group_name', 'api_name', 'method', 'url'], ['headers', 'params', 'body', 'type', 'updated_at']);
    await fetchHistory();
  } catch (e) {
    console.error(e)
    message.error(e.message);
  } finally {
    loading.value.request = false;
  }
};


const prepareSendRequest = async () => {
  let requestParams = {
    json: "",
    form: {},
    files: {},
  };
  let headers = header_editor.value.getValue();
  let query = query_editor.value.getValue();
  let body = body_editor.value.getValue();
  let headersObj = headers ? JSON.parse(headers) : {};
  let bodyObj = body ? JSON.parse(body) : {};
  let queryObj = query ? JSON.parse(query) : {};

  // 根据 bodyType 设置不同的请求体
  switch (bodyType.value) {
    case 'json':
      requestParams.json = body;
      break;
    case 'form':
      // 如果form body里的key以@ 开头，则认为是文件上传，真实key是去掉这部分的内容
      if (bodyObj) {
        bodyObj = JSON.parse(body)
        for (let key in bodyObj) {
          if (key.startsWith('@')) {
            requestParams.files[key.substring(1)] = bodyObj[key];
            delete bodyObj[key];
          }
        }
      }
      requestParams.form = bodyObj;
      break;
  }
  return {
    headersObj,
    queryObj,
    requestParams,
  };
}

// 删除历史记录
const deleteHistory = async (row) => {
  try {
    const sql = `DELETE FROM api WHERE id = '${row.id}'`;
    await ProxySql(sql);
    message.success('删除成功');
    await fetchHistory();
  } catch (e) {
    message.error('删除失败: ' + e.message);
    console.error(e.message)

  }
};

// 修改获取历史记录函数，支持搜索
const fetchHistory = async () => {
  loading.value.history = true;
  try {
    const offset = historyPagination.value.pageSize * (historyPagination.value.page - 1);
    let whereClause = '';
    const conditions = [];

    if (searchForm.value.groupName) {
      conditions.push(`group_name = '${searchForm.value.groupName}'`);
    }
    if (searchForm.value.apiName) {
      conditions.push(`api_name LIKE '%${searchForm.value.apiName}%'`);
    }
    if (searchForm.value.url) {
      conditions.push(`url LIKE '%${searchForm.value.url}%'`);
    }

    if (conditions.length > 0) {
      whereClause = 'WHERE ' + conditions.join(' AND ');
    }

    const sql = `
      SELECT *, count(*) over() AS total_count
      FROM api
      ${whereClause}
      ORDER BY updated_at DESC
      LIMIT ${historyPagination.value.pageSize} OFFSET ${offset}
    `;
    const history = await ProxyQuery(sql);
    historyList.value = history || [];
    if (history && history.length > 0) {
      historyPagination.value.itemCount = history[0].total_count;
    } else {
      historyPagination.value.itemCount = 0;
    }
  } catch (e) {
    message.error(e.message);
    console.error(e.message)

  } finally {
    loading.value.history = false;
  }
};

// 回填表单
const fillForm = (row) => {
  form.value = {
    groupName: row.group_name,
    apiName: row.api_name,
    method: row.method,
    url: row.url,
  };
  bodyType.value = row.type;

  header_editor.value.setValue(row.headers);
  query_editor.value.setValue(row.params);
  body_editor.value.setValue(row.body);

  message.success('已回填');
};

const clear = () => {
  header_editor.value.setValue(null);
  query_editor.value.setValue(null);
  body_editor.value.setValue(null);
  resp_editor.value.setValue(null);

  form.value = {
    groupName: '',
    apiName: '',
    method: 'GET',
    url: '',
  };
}


const benchmark = async () => {
  loading.value.benchmark = true;
  benchmarkForm.value.resp = '';
  try {
    let {
      headersObj,
      queryObj,
      requestParams,
    } = await prepareSendRequest();

    let res = await Benchmark(
        benchmarkForm.value.clientNum,
        benchmarkForm.value.secs,
        benchmarkForm.value.timeout,
        form.value.method,
        form.value.url,
        headersObj,
        queryObj,
        requestParams.json || "",
        requestParams.form || {},
        requestParams.files || {},
        false // useToken
    );
    benchmarkForm.value.resp = JSON.stringify(JSON.parse(res), null, 4);
  } catch (e) {
    message.error(e.message);
    console.error(e.message)
  } finally {
    loading.value.benchmark = false;
  }
};

// 初始化
onMounted(async () => {
  await Promise.allSettled(
      [
        fetchHistory(),
        fetchGroupOptions()
      ]
  )

});
</script>

<style scoped>
/* 根据需要添加样式 */
</style>