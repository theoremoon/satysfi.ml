module Main exposing (main)

import Browser
import Browser.Navigation
import Html exposing (..)
import Html.Attributes exposing (attribute, class, src, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode
import Url exposing (Url)



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = OnUrlRequest
        , onUrlChange = OnUrlChange
        }



-- MODEL


type FileTree
    = File String
    | Directory ( String, List FileTree, List FileTree )


type alias Model =
    { currentSource : String
    , product : Maybe String
    , fileTree : Maybe FileTree
    }


type Msg
    = OnUrlChange Url
    | OnUrlRequest Browser.UrlRequest
    | FileTreeRequest
    | FileTreeResult (Result Http.Error FileTree)
    | CompileRequest
    | CompileResult (Result Http.Error String)
    | SourceUpdate String



-- DECODER


makeDirectory : String -> List FileTree -> List FileTree -> FileTree
makeDirectory name dirs children =
    Directory ( name, dirs, children )


fileTreeFileDecoder : Json.Decode.Decoder FileTree
fileTreeFileDecoder =
    Json.Decode.map File (Json.Decode.field "name" Json.Decode.string)


fileTreeDecoder : Json.Decode.Decoder FileTree
fileTreeDecoder =
    Json.Decode.map3 makeDirectory
        (Json.Decode.field "name" Json.Decode.string)
        (Json.Decode.field "childdirs" (Json.Decode.list (Json.Decode.lazy (\_ -> fileTreeDecoder))))
        (Json.Decode.field "children" (Json.Decode.list fileTreeFileDecoder))



-- INIT


initialSource =
    """
@require: stdjabook

document (|
  title = {Hello, World!};
  author = {\\@theoremoon};
  show-toc = false;
  show-title = true;
|) '<
  +p {
    Hello, World!
  }
>
"""


init : () -> Url -> Browser.Navigation.Key -> ( Model, Cmd Msg )
init _ url navKey =
    ( Model initialSource Nothing Nothing
    , fileTreeRequest
    )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SourceUpdate newSource ->
            ( { model | currentSource = newSource }, Cmd.none )

        FileTreeRequest ->
            ( model, fileTreeRequest )

        FileTreeResult (Result.Ok fileTreeResult) ->
            ( { model | fileTree = Just fileTreeResult }, Cmd.none )

        FileTreeResult (Result.Err _) ->
            ( model, Cmd.none )

        CompileRequest ->
            ( model, compileRequest model.currentSource )

        CompileResult (Result.Ok result) ->
            ( { model | product = Just result }, Cmd.none )

        CompileResult (Result.Err _) ->
            ( model, Cmd.none )

        _ ->
            ( model, Cmd.none )



-- CMDS


fileTreeRequest : Cmd Msg
fileTreeRequest =
    Http.get
        { url = "/filetree"
        , expect = Http.expectJson FileTreeResult fileTreeDecoder
        }


compileRequest : String -> Cmd Msg
compileRequest source =
    Http.post
        { url = "/compile"
        , body = Http.stringBody "text/plain" source
        , expect = Http.expectString CompileResult
        }



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Hello World"
    , body =
        [ div
            [ class "container"
            ]
            [ div [ class "menu" ]
                [ input [ type_ "checkbox", class "fileTreeSwitch" ] [ text "FILES" ]
                , fileTree [ class "fileTree" ] model
                , button [ onClick CompileRequest ] [ text "COMPILE" ]
                ]
            , div [ class "main" ]
                [ satysfiEditor [] model
                , productViewer [] model
                ]
            ]
        ]
    }


satysfiEditor : List (Attribute Msg) -> Model -> Html Msg
satysfiEditor attrs model =
    textarea
        (attrs
            ++ [ class "editor"
               , value model.currentSource
               , onInput SourceUpdate
               ]
        )
        []


productViewer : List (Attribute Msg) -> Model -> Html Msg
productViewer attrs model =
    case model.product of
        Just pdf ->
            embed
                (attrs
                    ++ [ src ("data:application/pdf;base64," ++ pdf)
                       , class "viewer"
                       ]
                )
                []

        Nothing ->
            div attrs [ text "Nothing there" ]


fileTreeImpl : FileTree -> Html Msg
fileTreeImpl filetree =
    case filetree of
        File name ->
            li [] [ text name ]

        Directory ( name, dirs, children ) ->
            li []
                [ text name
                , ul [] (List.map fileTreeImpl dirs)
                , ul [] (List.map fileTreeImpl children)
                ]


fileTree : List (Attribute Msg) -> Model -> Html Msg
fileTree attrs model =
    case model.fileTree of
        Just tree ->
            ul attrs [ fileTreeImpl tree ]

        Nothing ->
            ul attrs []
