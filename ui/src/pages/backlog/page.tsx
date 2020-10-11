import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { cloneDeep, get, concat, isUndefined } from "lodash";
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import DeleteIcon from '@material-ui/icons/Delete';
import DoneIcon from '@material-ui/icons/Done';
import UndoIcon from '@material-ui/icons/Undo';
import { TaskAPI, Task } from "./../../services/tasks";
import { ListAPI, List } from "./../../services/list";
import { ProjectAPI } from "./../../services/project";
import { BacklogViewAPI, BacklogList, BacklogTask, BacklogProject } from "../../services/views/backlog";
import { FullScreenDialog } from './../../components/fullscreen-dialog';
import { Fab } from './../../components/fab';
import { Task as TaskComponent } from './../../components/task/task';
import { Table as TableComponent, RowMenuItem } from './../../components/table';


const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
  },
  container: {
    height: '100%',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  modal: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  paper: {
    backgroundColor: theme.palette.background.paper,
    border: '2px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
  },
}));


interface IProps {
  TaskAPI: TaskAPI;
  ListAPI: ListAPI;
  ProjectAPI: ProjectAPI;
  BacklogViewAPI: BacklogViewAPI;
}

interface Row extends BacklogTask { }

interface IState {
  lists: BacklogList[];
  projects: BacklogProject[];
}

interface SelectionChanged {
  id: string;
  name: string;
}

interface Selection {
  id: string;
  name: string;
}


function makeTaskCompletionButton(row: Row, onClick: (action: string) => void): RowMenuItem {
  const icon = row.status.completed ? UndoIcon : DoneIcon;
  const title = !row.status.completed ? 'Complete' : 'Reopen'
  return {
    title,
    icon,
    onClick: () => onClick(title),
  }
}

export function BacklogPage(props: IProps) {
  const classes = useStyles();
  const [selectedTaskIndex, updateSelectedTaskIndex] = useState(-1);
  const [taskEdit, updateTaskEdit] = useState<Task | undefined>();

  const [rows, updateRows] = useState<Row[]>([]);
  const [state, setState] = useState<IState>(() => {
    return {
      lists: [],
      projects: [],
    }
  });

  const [lists, updateLists] = useState<List[]>([])

  useEffect(() => {
    (async () => {
      const view = await props.BacklogViewAPI.get()
      const state: IState = {
        lists: view.lists,
        projects: view.projects,
      }
      setState(state)
      updateRows(view.tasks)

      const l = await props.ListAPI.list()
      updateLists(l);
    })()
  }, [props.ListAPI, props.TaskAPI, props.BacklogViewAPI])

  const onUpdate = async (newData: Row, oldData?: Row): Promise<any> => {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve();
        updateRows((prevState) => {
          let index;
          for (let i = 0; i < prevState.length; i++) {
            const element = prevState[i];
            if (element.metadata.id === newData.metadata.id) {
              index = i
            }
          }
          // index does not found, return previous state
          if (isUndefined(index)) {
            return prevState
          };
          const data = [...prevState];
          data[index] = newData;
          return data;
        });
      }, 600);
    });
  }

  const deleteTask = async (row: Row): Promise<void> => {
    await props.TaskAPI.remove(row.metadata.id)
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve();
        updateRows((prevState) => {
          let index;
          for (let i = 0; i < prevState.length; i++) {
            const element = prevState[i];
            if (element.metadata.id === row.metadata.id) {
              index = i
            }
          }
          // index does not found, return previous state
          if (isUndefined(index)) {
            return prevState
          };
          const data = [...prevState];
          data.splice(index, 1)
          return data;
        });
      }, 600);
    });
  }

  return (
    <Paper className={classes.root}>
      <TableComponent
        columns={[
          { title: "Title" },
          { title: "Description" },
          { title: "List" },
          { title: "Project" },
        ]}
        data={rows.map(r => {
          return {
            data: [
              {
                content: r.metadata.name,
              },
              {
                content: r.metadata.description,
              },
              {
                content: r.list.name,
              },
              {
                content: r.project.name,
              },
            ],
            onClick: (r: number) => {
              updateSelectedTaskIndex(r);
              updateTaskEdit(rows[r])
            },
            menu: [
              makeTaskCompletionButton(r, async (action: string) => {
                if (action === 'Reopen') {
                  await props.TaskAPI.reopen(r.metadata.id);
                }

                if (action === 'Complete') {
                  await props.TaskAPI.complete(r.metadata.id);
                }
                const newRow = cloneDeep(r)
                newRow.status.completed = !newRow.status.completed
                onUpdate(newRow, r)
              }),
              {
                icon: DeleteIcon,
                onClick: (i: number) => {
                  deleteTask(rows[i]);
                },
                title: 'Delete'
              }
            ],
          }
        })}
      />
      { selectedTaskIndex !== -1 &&
        <FullScreenDialog
          onClose={() => {
            updateSelectedTaskIndex(-1);
          }}
          component={
            <TaskComponent
              onChange={async () => {
                const view = await props.BacklogViewAPI.get()
                updateRows(view.tasks)
              }}
              defaultList={{
                id: rows[selectedTaskIndex].list.id,
                name: rows[selectedTaskIndex].list.name,
              }}
              defaultProject={{
                id: rows[selectedTaskIndex].project.id,
                name: rows[selectedTaskIndex].project.name,
              }}
              projects={state.projects.map(p => ({ name: p.name, id: p.id }))}
              lists={lists.map(l => ({ name: l.metadata.name, id: l.metadata.id }))}
              new={false}
              task={taskEdit as Task}
              TaskAPI={props.TaskAPI}
              ListAPI={props.ListAPI}
              ProjectAPI={props.ProjectAPI}
            />
          }
        />
      }
      <Fab
        includeFade={false}
        modal={
          <FullScreenDialog
          onClose={() => {
            updateSelectedTaskIndex(-1);
          }}
          component={
            <TaskComponent
              onChange={async () => {
                const view = await props.BacklogViewAPI.get()
                updateRows(view.tasks)
              }}
              defaultList={{
                id: get(rows, `[${selectedTaskIndex}].list.id`),
                name: get(rows, `[${selectedTaskIndex}].list.name`),
              }}
              defaultProject={{
                id: get(rows, `[${selectedTaskIndex}].project.id`),
                name: get(rows, `[${selectedTaskIndex}].project.name`),
              }}
              projects={state.projects.map(p => ({ name: p.name, id: p.id }))}
              lists={lists.map(l => ({ name: l.metadata.name, id: l.metadata.id }))}
              new={true}
              task={{ metadata: { id: '', name: '', description: '' }, spec: {}, status: { completed: false } }}
              TaskAPI={props.TaskAPI}
              ListAPI={props.ListAPI}
              ProjectAPI={props.ProjectAPI}
            />
          }
          />
        }
      />
    </Paper>
  )
}
