import React, { useState } from 'react';
import clsx from 'clsx';
import Grid from '@material-ui/core/Grid';
import { ProjectsPage } from './pages/projects/page'
import { ProjectPage } from './pages/project/page'
import { EmptyPage } from './pages/empty/page'
import { BacklogPage } from './pages/backlog/page'
import { HomePage } from './pages/home/page'
import { TriggersPage } from './pages/triggers/page'
import { LoginPage, getAPIToken } from './pages/login'
import { makeStyles } from '@material-ui/core/styles';
import AccountTreeIcon from '@material-ui/icons/AccountTree';
import RotateLeftIcon from '@material-ui/icons/RotateLeft';
import FlashAutoIcon from '@material-ui/icons/FlashAuto';
import FormatListBulletedIcon from '@material-ui/icons/FormatListBulleted';
import HomeIcon from '@material-ui/icons/Home';
import { Navbar } from './components/navbar/navbar'
import { Menu } from './components/menu/menu'
import { API as ProjectAPI } from './services/project'
import { API as TaskAPI } from './services/tasks'
import { API as ListAPI } from './services/list'
import { API as UserAPI } from './services/user'
import { API as BacklogViewAPI } from './services/views/backlog'
import { API as ProjectsViewAPI } from './services/views/projects'
import { API as ProjectViewAPI } from './services/views/project'
import { API as HomeViewAPI } from './services/views/home'
import { API as TriggersViewAPI } from './services/views/triggers'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect
} from "react-router-dom";
import './App.css';

const drawerWidth = 240;
const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  card: {
    background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
    border: 0,
    borderRadius: 3,
    boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
    color: 'white',
    height: 48,
    padding: '0 30px',
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
  appBar: {
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: drawerWidth,
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  hide: {
    display: 'none',
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(3),
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: -drawerWidth,
  },
  contentShift: {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: 0,
  },
}));



function App() {
  const classes = useStyles();
  const [menuOpen] = useState(false);
  const [APIToken, setAPIToken] = useState(getAPIToken());

  const handleMenuItemClicked = (name: string) => {
    console.log(`Item: ${name} clicked`);
  }
  if (APIToken === "") {
    return (
      <Router>
        <div className={classes.root}>
          <Switch>
            <Route path="/login" >
              <LoginPage UserAPI={UserAPI} onLogin={() => setAPIToken(getAPIToken())} />
            </Route>
          </Switch>
          <Redirect to="/login" />
        </div>
      </Router>
    )
  } else {
    UserAPI.updateDefaultsWithAPIToken(getAPIToken())
  }
  return (
    <Router>
      <div className={classes.root}>
        <Navbar className={clsx(classes.appBar, {
          [classes.appBarShift]: menuOpen,
        })} />
        <Menu
          drawerClassName={classes.drawer}
          drawerHeader={classes.drawerHeader}
          drawerOpen={true}
          drawerWidth={classes.drawerPaper}
          items={[
            {
              name: "home",
              text: 'Home',
              icon: HomeIcon,
              onClick: handleMenuItemClicked,
              link: "/"
            },
            {
              name: "projects",
              text: 'Projects',
              icon: AccountTreeIcon,
              onClick: handleMenuItemClicked,
              link: "/projects"
            },
            {
              name: "triggers",
              text: 'Triggers',
              icon: RotateLeftIcon,
              onClick: handleMenuItemClicked,
              link: "/triggers"
            },
            {
              name: "automation",
              text: 'Automation',
              icon: FlashAutoIcon,
              onClick: handleMenuItemClicked,
              link: "/automation"
            },
            {
              name: "backlog",
              text: 'Backlog',
              icon: FormatListBulletedIcon,
              onClick: handleMenuItemClicked,
              link: "/backlog"
            },
          ]} />
        <main
          className={clsx(classes.content, {
            [classes.contentShift]: true,
          })}
        >
          <div className={classes.drawerHeader} />
          <div>
            <Grid container spacing={3}>
              <Grid item xs={12}>

                <Switch>
                  {/* Order important */}
                  <Route path="/automation" >
                    <EmptyPage data={'automation'} />
                  </Route>
                  <Route path="/triggers">
                    <TriggersPage TriggersViewAPI={TriggersViewAPI} data={'triggers'} />
                  </Route>
                  <Route path="/backlog">
                    <BacklogPage ProjectAPI={ProjectAPI} TaskAPI={TaskAPI} ListAPI={ListAPI} BacklogViewAPI={BacklogViewAPI} />
                  </Route>
                  <Route path="/projects/:id" render={(rprops) => (<ProjectPage {...rprops} ListAPI={ListAPI} TaskAPI={TaskAPI} ProjectViewAPI={ProjectViewAPI} ProjectAPI={ProjectAPI}></ProjectPage>)}>
                  </Route>
                  <Route path="/projects">
                    <ProjectsPage ProjectsViewAPI={ProjectsViewAPI} ProjectAPI={ProjectAPI} />
                  </Route>
                  <Route path="/">
                    <HomePage HomeViewAPI={HomeViewAPI} ProjectAPI={ProjectAPI} TaskAPI={TaskAPI} ListAPI={ListAPI} />
                  </Route>
                </Switch>
              </Grid>
            </Grid>
          </div>
        </main>
      </div>
    </Router>
  );
}

export default App;
