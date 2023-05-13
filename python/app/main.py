import argparse
import os
import textwrap
import typing as t
from logging import getLogger

logger = getLogger(__name__)


def main(args: argparse.Namespace) -> None:
    """Example.

    >>> 1 + 1
    2
    """
    logger.info(args)
    return


if __name__ == "__main__":
    import logging
    from sys import stdout

    logging.basicConfig(
        stream=stdout,
        format="%(asctime)s %(levelname)s: %(filename)s %(lineno)d: %(message)s",
        level=logging.INFO,
    )

    h = """Sample command.
    """
    parser = argparse.ArgumentParser(
        description=textwrap.dedent(h).strip(),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    subparsers = parser.add_subparsers()

    options: dict[str, dict[str, t.Any]] = {
        "--stage": dict(
            required=True, default=os.environ.get("STAGE", "development")
        ),
        "--dt": dict(required=True, help="YYYY-MM-DD"),
    }

    def __register(handler: t.Callable, keys: list[str]) -> None:
        cmd = subparsers.add_parser(
            name=handler.__name__,
            description=handler.__doc__,
            help=handler.__doc__,
        )
        for key in keys:
            cmd.add_argument(key, **options[key])
        cmd.set_defaults(handler=handler)

    __register(main, ["--stage", "--dt"])

    args = parser.parse_args()
    if hasattr(args, "handler"):
        args.handler(args)
    else:
        parser.print_help()
