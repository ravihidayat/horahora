import { type IDownload } from "./types";
import { LinearProgress } from '@mui/material/LinearProgress';

import { blockComponent } from "#components/meta";
import {
  Card,
  CardBody,
  CardHeader,
  DL,
  DS,
  type ICardProps,
} from "#components/lists";
import { Heading } from "#components/headings";

export interface IDownloadCardProps extends ICardProps {
  download: IDownload;
}

export const DownloadCard = blockComponent(undefined, Component);

function Component({
  download,
  headingLevel,
  ...blockProps
}: IDownloadCardProps) {
  const { progress, DlStatus, VideoID, Website } = download;

  return (
    <Card {...blockProps}>
      <CardHeader>
        <Heading level={headingLevel}>Video {VideoID}</Heading>
      </CardHeader>
      <CardBody>
        <DL>
          <DS isHorizontal dKey="Website" dValue={Website} />
          <DS isHorizontal dKey="Status" dValue={DlStatus} />
          <DS
            isHorizontal
            dKey="Progress"
            dValue={<LinearProgress variant="determinate" value={progress} />}
          />
        </DL>
      </CardBody>
    </Card>
  );
}
