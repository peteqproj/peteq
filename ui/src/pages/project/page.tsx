import React, { useState , useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Typography from '@material-ui/core/Typography';
import { RouteComponentProps } from "react-router-dom";
import { ProjectAPI } from './../../services/project'
import { ProjectViewAPI, ProjectView } from './../../services/view.project'

const useStyles = makeStyles({
  root: {
    display: 'flex'
  },
  media: {
    height: 320,
    width: 320
  },
});


interface IProps extends RouteComponentProps {
  id?: string;
  ProjectAPI: ProjectAPI;
  ProjectViewAPI: ProjectViewAPI;
}

export function ProjectPage (props: IProps) {
  const classes = useStyles();
  const [state, setState] = useState({ metadata: { name: '', id: '', description: ''}, tasks: []} as ProjectView);
  const id = (props.match.params as any)['id'];
  useEffect(() => {    
    const fetch = async () => {
      const prj = await props.ProjectViewAPI.get(id)
      setState(prj)
    }
    fetch();
  
  }, [props.ProjectAPI, id]);
  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardContent>
          <Typography gutterBottom variant="h3" component="h2">
            {state.metadata.name}
          </Typography>
          <Typography variant="body1" color="textSecondary" component="p">
            {state.metadata.description}
          </Typography>
          <Typography gutterBottom variant="h6" component="h2">
            Tasks:
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
            {state.tasks.map((t) => (<div>{t.metadata.name}</div>))}
          </Typography>
        </CardContent>
      </CardActionArea>
      <CardMedia
          className={classes.media}
          image="https://images-na.ssl-images-amazon.com/images/I/41FH9qC4BrL._SX379_BO1,204,203,200_.jpg"
          title="Contemplative Reptile"
        />
    </Card>
  );
}

