<script>
  import {
    Header,

    SideNav,
    SideNavItems,

    Grid,
    Row,
    Column,

    SkipToContent,
    Content,
  } from "carbon-components-svelte";

  import { Router, Link, Route } from "svelte-routing";

  import { onMount } from 'svelte'

  import Messages from './components/Messages.svelte'
  import LeaderBoard from './pages/LeaderBoard.svelte'
  import { SSE } from './sse.js'

  let isSideNavOpen = false;

  let es = new SSE();

  onMount(async () => {
    es.dial()
  });
</script>
<Router>
  <Header persistentHamburgerMenu={true} company="insim.go" platformName="LiveMap"  bind:isSideNavOpen>
    <div slot="skip-to-content">
      <SkipToContent />
    </div>

    <SideNav bind:isOpen={isSideNavOpen}>
      <SideNavItems>
        <li class:bx--side-nav__item="{true}">
          <Link class="bx--side-nav__link" to="/"><span class:bx--side-nav__link-text="{true}">Leaderboard</span></Link>
        </li>
        <li class:bx--side-nav__item="{true}">
          <Link class="bx--side-nav__link" to="lap"><span class:bx--side-nav__link-text="{true}">Lap Chart</span></Link>
        </li>
      </SideNavItems>
    </SideNav>
  </Header>

  <Content>

  <Grid noGutter>
      <Route path="about">
        about
      </Route>

      <Route path="blog">
        blog
      </Route>
      <Route path="/">
        <LeaderBoard/>
      </Route>


      <Row>
      <Column>
      <Messages />
      </Column>
      </Row>

    </Grid>
  </Content>
</Router>
