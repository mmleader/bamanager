<script>
  import { onMount } from 'svelte';
  import BrowserCard from './components/BrowserCard.svelte';
  import BrowserModal from './components/BrowserModal.svelte';
  import BrowserTable from './components/BrowserTable.svelte';
  import ProxyManagerComp from './components/ProxyManager.svelte';
  import { 
    ListInstances, 
    AddInstance, 
    UpdateInstance, 
    DeleteInstance, 
    StartInstanceWithProxy, 
    StopInstance,
    GetConfig,
    SetMinimizeToTray,
    CheckProxyDirect,
    ListProxies,
    AddProxy,
    UpdateProxy,
    DeleteProxy
  } from '../wailsjs/go/main/App';

  let instances = [];
  let proxies = [];
  let showModal = false;
  let showProxyManager = false;

  let editingInstance = null;
  let minimizeToTray = false; // 默认值，之后通过 GetConfig 更新
  let viewMode = localStorage.getItem('viewMode') || 'card';
  let minimalMode = localStorage.getItem('minimalMode') === 'true';

  $: localStorage.setItem('viewMode', viewMode);
  $: localStorage.setItem('minimalMode', String(minimalMode));

  async function loadInstances() {
    try {
      instances = await ListInstances() || [];
    } catch (err) {
      console.error("加载列表失败:", err);
      alert("配置加载失败: " + err);
    }
  }
  async function handleSave(event) {
    const data = event.detail;
    try {
      if (data.id) {
        // 更新逻辑：需要构造一个完整的 BrowserInstance 对象
        const instToUpdate = instances.find(i => i.id === data.id);
        const updated = {
            ...instToUpdate,
            sortNum: data.sortNum,
            name: data.name,
            path: data.path,
            userDataDir: data.userDataDir,
            args: data.args,
            tags: data.tags,
            proxyId: data.proxyId || '',
            incognito: data.incognito || false
        };
        await UpdateInstance(updated);
      } else {
        const inst = await AddInstance(data.sortNum, data.name, data.path, data.userDataDir, data.args, data.tags);
        // 如果选了代理，需要更新实例的 proxyId
        if (data.proxyId && inst) {
          inst.proxyId = data.proxyId;
          await UpdateInstance(inst);
        }
      }
      showModal = false;
      await loadInstances();
      await loadProxies();
    } catch (err) {
      alert("保存失败: " + err);
    }
  }

  async function handleStart(id, proxyId = '') {
    try {
      await StartInstanceWithProxy(id, proxyId);
      await loadInstances();
    } catch (err) {
      alert("启动失败: " + err);
    }
  }

  async function handleStop(id) {
    try {
      await StopInstance(id);
      await loadInstances();
    } catch (err) {
      alert("停止失败: " + err);
    }
  }

  async function handleDelete(id) {
    if (!confirm("确定要删除此指纹配置吗？")) return;
    try {
      await DeleteInstance(id);
      await loadInstances();
    } catch (err) {
      alert("删除失败: " + err);
    }
  }

  async function handleCheckProxyDirect(proxyId, target) {
    const result = await CheckProxyDirect(proxyId, target);
    return result;
  }

  async function handleToggleIncognito(id, incognito) {
    try {
      const inst = instances.find(i => i.id === id);
      if (inst) {
        const updated = { ...inst, incognito };
        await UpdateInstance(updated);
        await loadInstances();
      }
    } catch (err) {
      alert("切换无痕模式失败: " + err);
    }
  }

  function openAddModal() {
    editingInstance = null;
    showModal = true;
  }

  function openEditModal(instance) {
    editingInstance = JSON.parse(JSON.stringify(instance));
    showModal = true;
  }

  function openDuplicateModal(instance) {
    const copy = JSON.parse(JSON.stringify(instance));
    copy.id = null; // Clear ID to treat as new
    copy.name = copy.name + " (副本)";
    editingInstance = copy;
    showModal = true;
  }

  // 代理管理相关函数
  async function loadProxies() {
    try {
      proxies = await ListProxies() || [];
    } catch (err) {
      console.error("加载代理列表失败:", err);
    }
  }

  async function handleProxySave(event) {
    const data = event.detail;
    try {
      if (data.id) {
        await UpdateProxy(data);
      } else {
        await AddProxy(data.name, data.protocol, data.host, data.port, data.username, data.password);
      }
      await loadProxies();
    } catch (err) {
      alert("保存代理失败: " + err);
    }
  }

  async function handleProxyDelete(event) {
    try {
      await DeleteProxy(event.detail.id);
      await loadProxies();
    } catch (err) {
      alert("删除代理失败: " + err);
    }
  }

  onMount(async () => {
    try {
      const config = await GetConfig();
      if (config) {
        minimizeToTray = config.minimize_to_tray;
      }
    } catch (e) {
      console.error("加载配置失败", e);
    }
    loadInstances();
    loadProxies();
    const interval = setInterval(loadInstances, 3000); // 定时同步状态
    return () => clearInterval(interval);
  });
</script>

<main class="container">
  <header>
    <h1>指纹浏览器管理器</h1>
    <div style="display: flex; align-items: center; gap: 1rem;">
      <div class="view-toggle">
        <button class:active={viewMode === 'card'} on:click={() => viewMode = 'card'} title="卡片视图">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
        </button>
        <button class:active={viewMode === 'table'} on:click={() => viewMode = 'table'} title="表格视图">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
        </button>
      </div>
      <label class="toggle-switch" title="精简模式 (只显示名字、标签、代理位置)">
        <input type="checkbox" bind:checked={minimalMode} />
        <span class="slider"></span>
        <span class="label-text">精简模式</span>
      </label>
      <label class="toggle-switch">
        <input type="checkbox" bind:checked={minimizeToTray} on:change={(e) => SetMinimizeToTray(e.target.checked)} />
        <span class="slider"></span>
        <span class="label-text">关闭最小化</span>
      </label>
      <button class="btn" style="border: 1px solid var(--border-color);" on:click={() => showProxyManager = true}>
        <svg style="margin-right: 0.5rem;" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
        代理管理
      </button>
      <button class="btn btn-primary" on:click={openAddModal}>
        <svg style="margin-right: 0.5rem;" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        添加实例
      </button>
    </div>
  </header>

  {#if viewMode === 'table'}
    <BrowserTable 
      {instances}
      {minimalMode}
      {proxies}
      onStart={handleStart}
      onStop={handleStop}
      onEdit={openEditModal}
      onDuplicate={openDuplicateModal}
      onDelete={handleDelete}
      onToggleIncognito={handleToggleIncognito}
    />
  {:else}
    <div class="grid" class:minimal={minimalMode}>
      {#each instances as instance (instance.id)}
        <BrowserCard 
          {instance} 
          {minimalMode}
          {proxies}
          onStart={handleStart}
          onStop={handleStop}
          onEdit={openEditModal}
          onDuplicate={openDuplicateModal}
          onDelete={handleDelete}
          onToggleIncognito={handleToggleIncognito}
        />
      {:else}
        <div style="grid-column: 1 / -1; text-align: center; padding: 4rem; background: white; border-radius: 0.75rem; border: 1px dashed var(--border-color); color: var(--text-muted);">
          还没有指纹实例，点击上方“添加实例”开始吧。
        </div>
      {/each}
    </div>
  {/if}

  <BrowserModal 
    show={showModal} 
    instance={editingInstance}
    {proxies}
    on:save={handleSave} 
    on:close={() => { showModal = false; loadInstances(); }} 
  />

  <ProxyManagerComp
    show={showProxyManager}
    {proxies}
    onCheckProxy={handleCheckProxyDirect}
    on:save={handleProxySave}
    on:delete={handleProxyDelete}
    on:close={() => { showProxyManager = false; loadProxies(); }}
  />
</main>

<style>
  /* 已经在 style.css 中定义了大部分基础样式 */
  .grid.minimal {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)) !important;
  }
  
  
  /* Toggle Switch Styles */
  .toggle-switch {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    user-select: none;
    font-size: 0.9rem;
    color: var(--text-main);
  }

  .toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
    position: absolute; /* Remove from flow */
  }

  .slider {
    position: relative;
    width: 44px;
    height: 24px;
    background-color: #cbd5e1; /* slate-300 */
    border-radius: 9999px;
    transition: .3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 20px;
    width: 20px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    transition: .3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 50%;
    box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
  }

  input:checked + .slider {
    background-color: var(--primary-color);
  }

  input:checked + .slider:before {
    transform: translateX(20px);
  }

  .label-text {
    font-weight: 500;
  }
  
  .view-toggle {
    display: flex;
    background: #e2e8f0;
    border-radius: 0.375rem;
    padding: 0.125rem;
  }
  
  .view-toggle button {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.25rem 0.5rem;
    border: none;
    background: transparent;
    cursor: pointer;
    border-radius: 0.25rem;
    color: var(--text-muted);
  }
  
  .view-toggle button.active {
    background: white;
    color: var(--primary-color);
    box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  }
</style>
